package cors

import (
	"os"
	"strings"

	gincors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Use attaches CORS middleware. If CORS_ALLOWED_ORIGINS is unset or empty,
// all origins are allowed (credentials disabled). Otherwise origins are a
// comma-separated list, e.g. "https://app.example.com,http://localhost:3000".
func Use(r *gin.Engine) {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	cfg := gincors.DefaultConfig()
	var origins []string
	for _, p := range strings.Split(raw, ",") {
		if o := strings.TrimSpace(p); o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		cfg.AllowAllOrigins = true
		cfg.AllowCredentials = false
	} else {
		cfg.AllowOrigins = origins
		cfg.AllowCredentials = true
	}
	cfg.AllowHeaders = []string{
		"Origin", "Content-Length", "Content-Type", "Authorization",
		"X-Admin-Key", "X-Internal-Token", "X-Webhook-Secret",
	}
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(gincors.New(cfg))
}
