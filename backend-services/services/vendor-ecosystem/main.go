// @title Vendor Ecosystem API
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
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"

	_ "github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/docs"
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
	dsn := config.MustGet("POSTGRES_VENDOR_DSN")

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

	r := api.NewRouter(api.RouterConfig{
		Store:          &repo.Store{DB: pool},
		Daily:          &repo.DailyVendorStore{DB: pool},
		InternalToken:  config.GetString("INTERNAL_AGGREGATE_TOKEN"),
		EnableSwagger:  config.GetBool("ENABLE_SWAGGER", true),
		Validator:      validator,
		PlatformFeeBPS: config.GetInt("PLATFORM_FEE_BPS", 500),
	})

	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}

	go func() {
		sugar.Infow("vendor-ecosystem listening", "addr", addr)
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
