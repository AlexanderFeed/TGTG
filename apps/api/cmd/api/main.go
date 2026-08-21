package main

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
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database startup failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	if cfg.MigrateOnStart {
		if err := database.Migrate(db); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}

	var sender mailer.Sender = mailer.DevelopmentSender{}
	if cfg.EmailDelivery == "smtp" {
		sender, err = mailer.NewSMTP(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
		if err != nil {
			log.Error("mailer setup failed", "error", err)
			os.Exit(1)
		}
	}

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           httpapi.New(cfg, store.New(db), sender, log),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		log.Info("api started", "config", cfg.String())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("api shutdown failed", "error", err)
		os.Exit(1)
	}
	log.Info("api stopped")
}
