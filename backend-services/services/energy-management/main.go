// @title Energy Management API
// @version 1.0
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

	"github.com/shellworlds/BRLBX4.0/backend-services/services/energy-management/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/energy-management/internal/repo"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/config"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/db"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/logger"
	mig "github.com/shellworlds/BRLBX4.0/backend-services/pkg/migrate"
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

	pgDSN := config.MustGet("POSTGRES_KITCHEN_DSN")
	tsDSN := config.MustGet("POSTGRES_TIMESCALE_DSN")

	root := config.GetString("SERVICE_ROOT")
	if root == "" {
		root = "."
	}
	pgMigDir := filepath.Join(root, "migrations", "postgres")
	tsMigDir := filepath.Join(root, "migrations", "timescale")

	if err := mig.Up(pgDSN, pgMigDir); err != nil {
		sugar.Fatalw("migrate postgres", "error", err)
	}
	if err := mig.Up(tsDSN, tsMigDir); err != nil {
		sugar.Fatalw("migrate timescale", "error", err)
	}

	pgPool, err := db.Connect(ctx, pgDSN)
	if err != nil {
		sugar.Fatal(err)
	}
	defer pgPool.Close()

	tsPool, err := db.Connect(ctx, tsDSN)
	if err != nil {
		sugar.Fatal(err)
	}
	defer tsPool.Close()

	var validator *auth.Validator
	if dom := config.GetString("AUTH0_DOMAIN"); dom != "" {
		validator = auth.NewValidator(auth.Config{
			Domain:   dom,
			Audience: config.GetString("AUTH0_AUDIENCE"),
		})
	}

	r := api.NewRouter(api.RouterConfig{
		Validator:    validator,
		AdminKey:     config.GetString("ADMIN_API_KEY"),
		IngestBearer: config.GetString("INGEST_BEARER_TOKEN"),
		Kitchens:     &repo.KitchenStore{Pool: pgPool},
		Readings:     &repo.ReadingStore{Pool: tsPool},
	})

	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}

	go func() {
		sugar.Infow("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		sugar.Warnw("shutdown", "error", err)
	}
}
