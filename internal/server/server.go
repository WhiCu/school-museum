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

	"github.com/WhiCu/school-museum/internal/config"
	"github.com/WhiCu/school-museum/internal/store"
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
	return a.srv.Shutdown(ctx)
}

func NewApp(cfg *config.Config, log *slog.Logger) *App {
	//TODO: remove

	// umami := telemetry.NewUmami("http://umami:3000", "7db11537-6def-4b16-a9c6-0ae33ac0641a")

	s := store.New()

	r := bunrouter.New(
		// bunrouter.Use(func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		// 	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// 		// umami.Track(req.Request, "API Call")
		// 		return next(w, req)
		// 	}
		// }),
		bunrouter.Use(bunrouterotel.NewMiddleware(
			bunrouterotel.WithClientIP(),
		)),
	)

	api := humabunrouter.New(r, huma.DefaultConfig("school-museum", "0.1.0"))
	pingHandler(api)

	// r.Mount("/museum", webmuseum.NewHandler(s, log.WithGroup("web-museum")))
	museum := huma.NewGroup(api, "/museum")
	webmuseum.RegisterHandlers(museum, s, log.WithGroup("web-museum"))

	// r.Mount("/admin", webadmin.NewHandler(s, log.WithGroup("web-admin")))
	admin := huma.NewGroup(api, "/admin")
	webadmin.RegisterHandlers(admin, s, log.WithGroup("web-admin"))

	return &App{
		srv: http.Server{
			Handler:      otelhttp.NewHandler(r, "school-museum-api"),
			Addr:         cfg.Server.ServerAddr(),
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		log: log,
	}
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
