package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

func TestPostVendorValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{
		Store:         &repo.Store{},
		Daily:         &repo.DailyVendorStore{},
		EnableSwagger: false,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportsVendorBadUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{
		Store:         &repo.Store{},
		Daily:         &repo.DailyVendorStore{},
		EnableSwagger: false,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/vendor/not-uuid?from=2024-01-01&to=2024-01-02", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInternalAggregateDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{
		Store:         &repo.Store{},
		Daily:         &repo.DailyVendorStore{},
		InternalToken: "",
		EnableSwagger: false,
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/internal/aggregate/daily", nil))
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTransactionsBindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{Store: &repo.Store{}, Daily: &repo.DailyVendorStore{}, EnableSwagger: false})
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{"vendor_id": "x"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
