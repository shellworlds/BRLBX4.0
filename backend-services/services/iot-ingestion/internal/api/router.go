package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/repo"
)

type RouterConfig struct {
	Store         *repo.Store
	EnableSwagger bool
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", metrics.Handler())
	if cfg.EnableSwagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := r.Group("/api/v1")
	v1.GET("/alerts", func(c *gin.Context) {
		items, err := cfg.Store.ListAlerts(c.Request.Context(), 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})
	v1.POST("/alerts/ack", func(c *gin.Context) {
		var body struct {
			ID uuid.UUID `json:"id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := cfg.Store.AckAlert(c.Request.Context(), body.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}
