package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/auth-rbac/internal/repo"
)

func TestWebhookDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(RouterConfig{
		Validator:     auth.NewValidator(auth.Config{Domain: "example.com", Audience: "a"}),
		WebhookSecret: "",
		Store:         &repo.Store{},
		EnableSwagger: false,
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/users/sync", bytes.NewReader([]byte("{}"))))
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
}
