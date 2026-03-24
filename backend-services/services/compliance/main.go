package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/config"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/db"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/logger"
	mig "github.com/shellworlds/BRLBX4.0/backend-services/pkg/migrate"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/compliance/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/compliance/internal/repo"
)

func main() {
	config.Load()
	zlog, err := logger.New(config.GetBool("LOG_DEV", false))
	if err != nil {
		panic(err)
	}
	defer zlog.Sync() //nolint:errcheck
	sugar := zlog.Sugar()

	ctx := context.Background()
	dsn := config.MustGet("POSTGRES_COMPLIANCE_DSN")

	root := config.GetString("SERVICE_ROOT")
	if root == "" {
		root = "."
	}
	if err := mig.Up(dsn, filepath.Join(root, "migrations", "postgres")); err != nil {
		sugar.Fatalw("migrate", "error", err)
	}

	pool, err := db.Connect(ctx, dsn)
	if err != nil {
		sugar.Fatal(err)
	}
	defer pool.Close()

	var validator *auth.Validator
	if dom := config.GetString("AUTH0_DOMAIN"); dom != "" {
		validator = auth.NewValidator(auth.Config{
			Domain:   dom,
			Audience: config.GetString("AUTH0_AUDIENCE"),
		})
	}

	sales := strings.TrimSpace(config.GetString("SALES_NOTIFY_EMAIL"))
	if sales == "" {
		sales = "sales@borelsigma.com"
	}

	r := api.NewRouter(api.RouterConfig{
		Store:           &repo.Store{DB: pool},
		Validator:       validator,
		SMTPHost:        config.GetString("SMTP_HOST"),
		SMTPPort:        config.GetString("SMTP_PORT"),
		SMTPUser:        config.GetString("SMTP_USER"),
		SMTPPass:        config.GetString("SMTP_PASS"),
		MailFrom:        config.GetString("SMTP_FROM"),
		SalesRecipients: sales,
		EnableSwagger:   config.GetBool("ENABLE_SWAGGER", false),
	})

	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}

	go func() {
		sugar.Infow("compliance listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
