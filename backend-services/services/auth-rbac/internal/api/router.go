package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/auth-rbac/internal/repo"
)

type RouterConfig struct {
	Validator     *auth.Validator
	WebhookSecret string
	Store         *repo.Store
	EnableSwagger bool
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	cors.Use(r)
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", metrics.Handler())
	if cfg.EnableSwagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := r.Group("/api/v1")

	v1.POST("/users/sync", func(c *gin.Context) {
		if cfg.WebhookSecret == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "webhook not configured"})
			return
		}
		provided := c.GetHeader("X-Webhook-Secret")
		if provided == "" {
			h := c.GetHeader("Authorization")
			if strings.HasPrefix(strings.ToLower(h), "bearer ") {
				provided = strings.TrimSpace(h[7:])
			}
		}
		if provided != cfg.WebhookSecret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook secret"})
			return
		}
		var body struct {
			Auth0ID  string  `json:"auth0_id"`
			Sub      string  `json:"sub"`
			Email    string  `json:"email" binding:"required"`
			Role     string  `json:"role" binding:"required"`
			ClientID *string `json:"client_id"`
			VendorID *string `json:"vendor_id"`
			Region   string  `json:"region"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		auth0ID := strings.TrimSpace(body.Auth0ID)
		if auth0ID == "" {
			auth0ID = strings.TrimSpace(body.Sub)
		}
		if auth0ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "auth0_id or sub required"})
			return
		}
		u := &repo.User{Auth0ID: auth0ID, Email: body.Email, Role: strings.ToLower(body.Role), Region: strings.TrimSpace(body.Region)}
		u.ClientID = emptyStringPtr(body.ClientID)
		u.VendorID = emptyStringPtr(body.VendorID)

		if err := cfg.Store.Upsert(c.Request.Context(), u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	me := v1.Group("/users")
	me.Use(auth.Middleware(cfg.Validator))
	me.GET("/me", func(c *gin.Context) {
		sub := auth.SubjectFromContext(c)
		u, err := cfg.Store.GetByAuth0(c.Request.Context(), sub)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not synced"})
			return
		}
		perms := permissionsForRole(u.Role)
		c.JSON(http.StatusOK, gin.H{"user": u, "permissions": perms, "token_roles": auth.RolesFromContext(c)})
	})

	return r
}

func emptyStringPtr(p *string) *string {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return nil
	}
	return &s
}

func permissionsForRole(role string) []string {
	switch strings.ToLower(role) {
	case "admin":
		return []string{"*"}
	case "vendor":
		return []string{"vendor.read", "vendor.write", "kitchen.read"}
	case "client":
		return []string{"client.read"}
	default:
		return []string{}
	}
}
