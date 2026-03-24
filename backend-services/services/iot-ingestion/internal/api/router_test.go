package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/repo"
)

func TestAlertsRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	st := &repo.Store{} // nil PG — list will error
	r := NewRouter(RouterConfig{Store: st, EnableSwagger: false})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/alerts", nil))
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
