package server

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"strings"
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

func (a *App) shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) *App {
	app := &App{
		log:             log,
		shutdownTimeout: cfg.Server.ShutdownTimeout,
	}

	// ----- Database -----
	database, err := db.NewDB(ctx, cfg.Storage.DSN(), db.WithDebug(true))
	if err != nil {
		log.Error("failed to create database connection", slog.String("error", err.Error()))
		panic(err)
	}

	news := storage.NewNewsStorage(database)
	exhibits := storage.NewExhibitStorage(database)
	exhibitions := storage.NewExhibitionStorage(database)
	visits := storage.NewVisitStorage(database)

	// ----- Router -----
	r := bunrouter.New()

	api := humabunrouter.New(r, huma.DefaultConfig("school-museum", "0.1.0"))
	pingHandler(api)

	museum := huma.NewGroup(api, "/museum")
	webmuseum.RegisterHandlers(
		museum, news, exhibitions, exhibits, visits, log.WithGroup("web-museum"))

	admin := huma.NewGroup(api, "/admin")
	webadmin.RegisterHandlers(
		admin, news, exhibitions, exhibits, visits, log.WithGroup("web-admin"))

	// ----- HTTP handler chain -----
	var handler http.Handler = r

	// Basic Auth middleware for /admin/ routes
	handler = adminAuthMiddleware(handler, cfg.Admin.Login, cfg.Admin.Password)

	// Visit tracking middleware — extracts IP and User-Agent into request context
	handler = visitTrackingMiddleware(handler)

	// CORS middleware
	handler = corsMiddleware(handler)

	app.srv = http.Server{
		Handler:      handler,
		Addr:         cfg.Server.ServerAddr(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return app
}

// visitTrackingMiddleware extracts the visitor's IP address and User-Agent
// from the HTTP request and stores them in the request context.
// Downstream handlers can read them via model.CtxKeyVisitorIP / model.CtxKeyVisitorUA.
func visitTrackingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		ua := r.Header.Get("User-Agent")

		ctx := context.WithValue(r.Context(), model.CtxKeyVisitorIP, ip)
		ctx = context.WithValue(ctx, model.CtxKeyVisitorUA, ua)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractIP retrieves the real client IP, checking proxy headers first.
func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
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

// adminAuthMiddleware protects /admin/ routes with Basic Auth.
// Requests to other paths pass through unchanged.
func adminAuthMiddleware(next http.Handler, login, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin/") && r.URL.Path != "/admin" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, `{"title":"Unauthorized","status":401}`, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			http.Error(w, `{"title":"Unauthorized","status":401}`, http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			http.Error(w, `{"title":"Unauthorized","status":401}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			http.Error(w, `{"title":"Unauthorized","status":401}`, http.StatusUnauthorized)
			return
		}

		loginOk := subtle.ConstantTimeCompare([]byte(parts[0]), []byte(login)) == 1
		passOk := subtle.ConstantTimeCompare([]byte(parts[1]), []byte(password)) == 1

		if !loginOk || !passOk {
			http.Error(w, `{"title":"Unauthorized","status":401}`, http.StatusUnauthorized)
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
