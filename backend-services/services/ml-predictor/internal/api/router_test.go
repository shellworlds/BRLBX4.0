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

	"github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/internal/cache"
)

func TestModelsRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{Cache: cache.Noop{}, EnableSwagger: false})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/predict/models", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestEnergyPredictRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{Cache: cache.Noop{}, EnableSwagger: false})
	body, _ := json.Marshal(map[string]any{
		"kitchen_id":   uuid.New().String(),
		"hours_ahead":  2,
		"history_kw": []float64{5, 6, 7},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/predict/energy", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
