// @title IoT Ingestion API
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

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/config"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/db"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/logger"
	mig "github.com/shellworlds/BRLBX4.0/backend-services/pkg/migrate"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/anomaly"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/api"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/energy"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/repo"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/watchdog"

	_ "github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/docs"
)

func main() {
	config.Load()
	zlog, err := logger.New(config.GetBool("LOG_DEV", false))
	if err != nil {
		panic(err)
	}
	defer zlog.Sync() //nolint:errcheck
	sugar := zlog.Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgDSN := config.MustGet("POSTGRES_IOT_DSN")
	tsDSN := config.MustGet("POSTGRES_IOT_TIMESCALE_DSN")

	root := config.GetString("SERVICE_ROOT")
	if root == "" {
		root = "."
	}
	if err := mig.Up(pgDSN, filepath.Join(root, "migrations", "postgres")); err != nil {
		sugar.Fatalw("migrate postgres", "error", err)
	}
	if err := mig.Up(tsDSN, filepath.Join(root, "migrations", "timescale")); err != nil {
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

	st := &repo.Store{PG: pgPool, TS: tsPool}
	ec := &energy.Client{
		BaseURL:     config.MustGet("ENERGY_SERVICE_URL"),
		IngestToken: config.GetString("INGEST_BEARER_TOKEN"),
	}
	tr := watchdog.NewTracker()

	broker := config.MustGet("MQTT_BROKER_URL")
	clientID := config.GetString("MQTT_CLIENT_ID")
	if clientID == "" {
		clientID = "iot-ingestion"
	}
	topic := config.GetString("MQTT_TOPIC_FILTER")
	if topic == "" {
		topic = "borelsigma/kitchen/+/telemetry"
	}
	user := config.GetString("MQTT_USERNAME")
	pass := config.GetString("MQTT_PASSWORD")

	anom := anomaly.New(config.GetInt("ANOMALY_WINDOW", 20), config.GetFloat64("ANOMALY_SIGMA", 3))
	if config.GetBool("ANOMALY_DISABLED", false) {
		anom = nil
	}
	if err := startMQTT(ctx, broker, clientID, user, pass, topic, st, ec, tr, config.GetString("SLACK_WEBHOOK_URL"), anom); err != nil {
		sugar.Fatalw("mqtt", "error", err)
	}

	var caCert, caKey []byte
	if p := config.GetString("IOT_DEVICE_CA_CERT_FILE"); p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			sugar.Fatalw("read IOT_DEVICE_CA_CERT_FILE", "error", err)
		}
		caCert = b
	}
	if p := config.GetString("IOT_DEVICE_CA_KEY_FILE"); p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			sugar.Fatalw("read IOT_DEVICE_CA_KEY_FILE", "error", err)
		}
		caKey = b
	}
	if (len(caCert) > 0) != (len(caKey) > 0) {
		sugar.Fatalw("IOT_DEVICE_CA_CERT_FILE and IOT_DEVICE_CA_KEY_FILE must both be set or both omitted")
	}

	r := api.NewRouter(api.RouterConfig{
		Store:               st,
		EnableSwagger:       config.GetBool("ENABLE_SWAGGER", true),
		InternalDeviceToken: config.GetString("INTERNAL_DEVICE_TOKEN"),
		DeviceCACert:        caCert,
		DeviceCAKey:         caKey,
	})
	addr := fmt.Sprintf(":%d", config.GetInt("PORT", 8080))
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}

	go func() {
		sugar.Infow("iot-ingestion listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	cancel()
}
