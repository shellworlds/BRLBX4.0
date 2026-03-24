// @title Auth RBAC API
// @version 1.0
// @BasePath /
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/config"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/db"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/logger"
	mig "github.com/shellworlds/BRLBX4.0/backend-services/pkg/migrate"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/auth-rbac/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/auth-rbac/internal/repo"

	_ "github.com/shellworlds/BRLBX4.0/backend-services/services/auth-rbac/docs"
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
	dsn := config.MustGet("POSTGRES_AUTH_DSN")

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

	domain := config.MustGet("AUTH0_DOMAIN")
	validator := auth.NewValidator(auth.Config{
		Domain:   domain,
		Audience: config.GetString("AUTH0_AUDIENCE"),
	})

	r := api.NewRouter(api.RouterConfig{
		Validator:     validator,
		WebhookSecret: config.GetString("AUTH0_SYNC_WEBHOOK_SECRET"),
		Store:         &repo.Store{DB: pool},
		EnableSwagger: config.GetBool("ENABLE_SWAGGER", true),
	})

	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}

	go func() {
		sugar.Infow("auth-rbac listening", "addr", addr)
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
