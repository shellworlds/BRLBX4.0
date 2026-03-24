package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/energy"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/repo"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/watchdog"
)

func startMQTT(ctx context.Context, broker, clientID, user, pass, topic string, st *repo.Store, ec *energy.Client, tr *watchdog.Tracker, slackURL string) error {
	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetOrderMatters(false)
	if user != "" {
		opts.SetUsername(user).SetPassword(pass)
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	handler := func(_ MQTT.Client, msg MQTT.Message) {
		kitchenID, err := kitchenFromTopic(msg.Topic())
		if err != nil {
			log.Printf("iot: topic parse: %v", err)
			return
		}
		var reading energy.Reading
		if err := json.Unmarshal(msg.Payload(), &reading); err != nil {
			log.Printf("iot: json: %v", err)
			return
		}
		payload := json.RawMessage(msg.Payload())
		if err := st.InsertRaw(context.Background(), kitchenID, msg.Topic(), payload); err != nil {
			log.Printf("iot: raw insert: %v", err)
		}
		if err := ec.PostReading(context.Background(), kitchenID, reading); err != nil {
			log.Printf("iot: energy forward: %v", err)
		}
		tr.Seen(kitchenID, time.Now().UTC())
	}

	if token := client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				client.Disconnect(250)
				return
			case <-t.C:
				_ = tr.Check(context.Background(), 5*time.Minute, func(ctx context.Context, kitchen uuid.UUID, last time.Time) error {
					msg := fmt.Sprintf("no telemetry for kitchen %s since %s", kitchen, last.Format(time.RFC3339))
					if err := st.InsertAlert(ctx, kitchen, "offline", msg); err != nil {
						return err
					}
					postSlack(slackURL, msg)
					return nil
				})
			}
		}
	}()

	return nil
}

func kitchenFromTopic(topic string) (uuid.UUID, error) {
	// borelsigma/kitchen/{id}/telemetry
	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return uuid.Nil, fmt.Errorf("unexpected topic")
	}
	return uuid.Parse(parts[2])
}

func postSlack(url, text string) {
	if url == "" {
		return
	}
	body := map[string]any{"text": text}
	b, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := newHTTPPost(ctx, url, b)
	if err != nil {
		return
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
}

func newHTTPPost(ctx context.Context, url string, b []byte) (*http.Request, error) {
	// small indirection to allow testing with net/http
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

var httpClient = getHTTPClient()

func getHTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
