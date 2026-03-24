package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/energy-management/internal/repo"
)

func TestPostReadingIngestBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	readings := &repo.ReadingStore{} // nil pool exercises handler after auth check
	kitchens := &repo.KitchenStore{}
	cfg := RouterConfig{
		IngestBearer: "secret",
		Kitchens:     kitchens,
		Readings:     readings,
	}
	r := NewRouter(cfg)

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	body, err := json.Marshal(map[string]any{
		"grid_power": 1.0, "battery_power": 2.0, "solar_power": 3.0, "lpg_status": "ok", "uptime_percent": 99.0,
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/kitchens/"+id.String()+"/readings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer wrong")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/kitchens/"+id.String()+"/readings", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer secret")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	require.NotEqual(t, http.StatusUnauthorized, w2.Code)
}
