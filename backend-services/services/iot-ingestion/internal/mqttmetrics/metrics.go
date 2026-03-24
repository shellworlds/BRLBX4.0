package mqttmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	IngestMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_ingest_messages_total",
		Help: "MQTT telemetry messages processed by iot-ingestion",
	})
	AnomalyAlerts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "iot_anomaly_alerts_total",
		Help: "Anomaly detections that created ingestion alerts",
	})
)
