package mqttmetrics

import "testing"

func TestCounters(t *testing.T) {
	IngestMessages.Inc()
	AnomalyAlerts.Inc()
}
