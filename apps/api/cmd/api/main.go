package main

// This file is the backend executable's entry point.
//
// A Go program starts in package main, inside func main. main does not contain
// business rules; its job is to build the application by connecting the
// configuration, database, mail sender, HTTP routes, and operating-system
// shutdown signal. This style is often called the "composition root".

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eshche-est/eshche-est/apps/api/internal/config"
	"github.com/eshche-est/eshche-est/apps/api/internal/database"
	"github.com/eshche-est/eshche-est/apps/api/internal/httpapi"
	"github.com/eshche-est/eshche-est/apps/api/internal/mailer"
	"github.com/eshche-est/eshche-est/apps/api/internal/store"
)

func main() {
	// slog writes one JSON object per log line. Structured logs are easier for
	// Docker and log collectors to search than free-form fmt.Println output.
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load reads environment variables once at startup. If required settings
	// are missing, stopping immediately is safer than running a half-configured
	// server and discovering the problem during a user request.
	cfg, err := config.Load()
	if err != nil {
		log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	// The root context is cancelled by Ctrl+C (SIGINT) or by Docker/systemd
	// asking the process to stop (SIGTERM). Request contexts ultimately inherit
	// from the HTTP server, while this one controls application lifetime.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// sql.DB is a concurrency-safe connection pool, not a single connection.
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database startup failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	// Migrations make the database schema match the version expected by this
	// binary. The SQL files are embedded into the compiled Go executable.
	if cfg.MigrateOnStart {
		if err := database.Migrate(db); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}

	// Both senders satisfy mailer.Sender. The handlers therefore depend on a
	// small interface and do not need to know whether SMTP is actually used.
	var sender mailer.Sender = mailer.DevelopmentSender{}
	if cfg.EmailDelivery == "smtp" {
		sender, err = mailer.NewSMTP(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
		if err != nil {
			log.Error("mailer setup failed", "error", err)
			os.Exit(1)
		}
	}

	// This is where all backend layers are wired together:
	//
	//   PostgreSQL -> store.New -> httpapi.New -> net/http server
	//
	// httpapi.New returns an http.Handler (the router). net/http only knows how
	// to pass requests to that interface; it knows nothing about our features.
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           httpapi.New(cfg, store.New(db), sender, log),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	// ListenAndServe blocks, so it runs in a goroutine. main remains available
	// below to wait for the cancellation signal and perform a graceful shutdown.
	go func() {
		log.Info("api started", "config", cfg.String())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Receiving from Done blocks here until the process is asked to stop.
	<-ctx.Done()

	// Graceful shutdown stops accepting new requests and gives in-flight
	// requests a short deadline to finish before the process exits.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("api shutdown failed", "error", err)
		os.Exit(1)
	}
	log.Info("api stopped")
}
