package cors

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUse_DefaultAllowAll(t *testing.T) {
	_ = os.Unsetenv("CORS_ALLOWED_ORIGINS")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	Use(r)
	require.NotNil(t, r)
}

func TestUse_WithOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	Use(r)
	require.NotNil(t, r)
}
