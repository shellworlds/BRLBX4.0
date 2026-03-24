// @title ML Predictor API
// @version 1.0
// @BasePath /
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/config"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/logger"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/internal/cache"

	_ "github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/docs"
)

func main() {
	config.Load()
	zlog, err := logger.New(config.GetBool("LOG_DEV", false))
	if err != nil {
		panic(err)
	}
	defer zlog.Sync() //nolint:errcheck
	sugar := zlog.Sugar()

	var c cache.PredictCache = cache.Noop{}
	if addr := config.GetString("REDIS_ADDR"); addr != "" {
		c = cache.NewRedis(addr, config.GetString("REDIS_PASSWORD"), config.GetInt("REDIS_DB", 0), config.GetString("REDIS_PREFIX")+"ml:")
	}

	r := api.NewRouter(api.RouterConfig{
		Cache:         c,
		EnableSwagger: config.GetBool("ENABLE_SWAGGER", true),
	})

	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}

	go func() {
		sugar.Infow("ml-predictor listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
