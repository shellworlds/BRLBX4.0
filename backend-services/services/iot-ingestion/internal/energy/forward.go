package energy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Reading struct {
	Timestamp     *time.Time `json:"timestamp"`
	GridPower     float64    `json:"grid_power"`
	BatteryPower  float64    `json:"battery_power"`
	SolarPower    float64    `json:"solar_power"`
	LPGStatus     string     `json:"lpg_status"`
	UptimePercent float64    `json:"uptime_percent"`
}

type Client struct {
	BaseURL    string
	HTTP       *http.Client
	IngestToken string
}

func (c *Client) PostReading(ctx context.Context, kitchen uuid.UUID, r Reading) error {
	if c.HTTP == nil {
		c.HTTP = http.DefaultClient
	}
	url := fmt.Sprintf("%s/api/v1/kitchens/%s/readings", c.BaseURL, kitchen.String())
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.IngestToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.IngestToken)
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("energy: unexpected status %s", res.Status)
	}
	return nil
}
