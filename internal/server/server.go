package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/WhiCu/school-museum/db"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/WhiCu/school-museum/internal/config"
	"github.com/WhiCu/school-museum/internal/telemetry"
	webadmin "github.com/WhiCu/school-museum/internal/web-admin"
	webmuseum "github.com/WhiCu/school-museum/internal/web-museum"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humabunrouter"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
)

type App struct {
	srv http.Server

	log *slog.Logger

	shutdownTimeout time.Duration
	telemetry       *telemetry.Provider
}

func (a *App) gracefulShutdownCtx(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	a.log.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	if err := a.shutdown(ctx); err != nil {
		a.log.Error("server forced to shutdown with error", slog.String("error", err.Error()))
		return err
	}

	a.log.Info("server successfully shutdown")

	return nil
}

func (a *App) shutdown(ctx context.Context) (err error) {
	if a.telemetry != nil {
		if tErr := a.telemetry.Shutdown(ctx); tErr != nil {
			a.log.Error("telemetry shutdown error", slog.String("error", tErr.Error()))
		}
	}
	return a.srv.Shutdown(ctx)
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) *App {
	app := &App{
		log:             log,
		shutdownTimeout: cfg.Server.ShutdownTimeout,
	}

	// ----- Telemetry (OTLP) -----
	if cfg.Telemetry.Enabled && cfg.Telemetry.OTLPEndpoint != "" {
		tp, err := telemetry.InitOTLP(ctx, cfg.Telemetry.OTLPEndpoint, cfg.Telemetry.ServiceName, log.WithGroup("telemetry"))
		if err != nil {
			log.Error("failed to init OTLP telemetry, continuing without it", slog.String("error", err.Error()))
		} else {
			app.telemetry = tp
		}
	}

	// ----- Database -----
	dbOpts := []db.Option{db.WithDebug(true)}
	if app.telemetry != nil {
		dbOpts = append(dbOpts, db.WithOtel())
	}

	database, err := db.NewDB(ctx, cfg.Storage.DSN(), dbOpts...)
	if err != nil {
		log.Error("failed to create database connection", slog.String("error", err.Error()))
		panic(err)
	}

	news := storage.NewNewsStorage(database)
	exhibits := storage.NewExhibitStorage(database)
	exhibitions := storage.NewExhibitionStorage(database)

	// ----- Router -----
	r := bunrouter.New(
		bunrouter.Use(bunrouterotel.NewMiddleware(
			bunrouterotel.WithClientIP(),
		)),
	)

	api := humabunrouter.New(r, huma.DefaultConfig("school-museum", "0.1.0"))
	pingHandler(api)

	museum := huma.NewGroup(api, "/museum")
	webmuseum.RegisterHandlers(
		museum, news, exhibitions, exhibits, log.WithGroup("web-museum"))

	admin := huma.NewGroup(api, "/admin")
	webadmin.RegisterHandlers(
		admin, news, exhibitions, exhibits, log.WithGroup("web-admin"))

	// ----- HTTP handler chain -----
	var handler http.Handler = otelhttp.NewHandler(r, "school-museum-api")

	// CORS middleware
	handler = corsMiddleware(handler)

	// Umami analytics — resolve website at startup, expose config to frontend
	if cfg.Telemetry.Umami.Enabled && cfg.Telemetry.Umami.URL != "" {
		umamiInfo, err := telemetry.ResolveUmamiWebsite(ctx, telemetry.UmamiOpts{
			URL:       cfg.Telemetry.Umami.URL,
			WebsiteID: cfg.Telemetry.Umami.WebsiteID,
			Username:  cfg.Telemetry.Umami.Username,
			Password:  cfg.Telemetry.Umami.Password,
			Domain:    cfg.Telemetry.Umami.Domain,
		}, log.WithGroup("umami"))
		if err != nil {
			log.Error("failed to resolve umami website, analytics disabled", slog.String("error", err.Error()))
		} else {
			analyticsHandler(api, umamiInfo)
			log.Info("umami analytics enabled (frontend tracking)",
				slog.String("url", umamiInfo.URL),
				slog.String("website_id", umamiInfo.WebsiteID))
		}
	}

	app.srv = http.Server{
		Handler:      handler,
		Addr:         cfg.Server.ServerAddr(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return app
}

// corsMiddleware adds CORS headers to allow cross-origin requests from the frontend.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		a.log.Info("starting http server", slog.String("addr", a.srv.Addr))
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server failed", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	eg.Go(func() error {
		return a.gracefulShutdownCtx(ctx)
	})

	fmt.Println(`
	==================================
	=                                =
	=    Server successfully start   =
	=                                =
	==================================
	`)

	return eg.Wait()
}
