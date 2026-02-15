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
	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/WhiCu/school-museum/internal/config"
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

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) *App {
	//TODO: remove

	db, err := db.NewDB(ctx, cfg.Storage.DSN())
	if err != nil {
		log.Error("failed to create database connection", slog.String("error", err.Error()))
		panic(err)
	}

	news := storage.NewNewsStorage(db)
	exhibits := storage.NewExhibitStorage(db)
	exhibitions := storage.Storage[model.Exhibition](nil)

	r := bunrouter.New(
		bunrouter.Use(bunrouterotel.NewMiddleware(
			bunrouterotel.WithClientIP(),
		)),
	)

	api := humabunrouter.New(r, huma.DefaultConfig("school-museum", "0.1.0"))
	pingHandler(api)

	// r.Mount("/museum", webmuseum.NewHandler(s, log.WithGroup("web-museum")))
	museum := huma.NewGroup(api, "/museum")
	webmuseum.RegisterHandlers(
		museum, news, exhibitions, exhibits, log.WithGroup("web-museum"))

	// r.Mount("/admin", webadmin.NewHandler(s, log.WithGroup("web-admin")))
	admin := huma.NewGroup(api, "/admin")
	webadmin.RegisterHandlers(
		admin, news, exhibitions, exhibits, log.WithGroup("web-admin"))

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
