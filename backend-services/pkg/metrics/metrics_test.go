package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/m", Handler())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/m", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
