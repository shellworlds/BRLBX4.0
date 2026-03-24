package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/devcert"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/iot-ingestion/internal/repo"
)

type RouterConfig struct {
	Store         *repo.Store
	EnableSwagger bool
	// InternalDeviceToken protects POST /api/v1/internal/devices/register.
	InternalDeviceToken string
	// Device CA material (PEM). When both are set, CSRs are signed for MQTT mTLS client certs.
	DeviceCACert []byte
	DeviceCAKey  []byte
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

	if cfg.InternalDeviceToken != "" && cfg.Store != nil {
		v1.POST("/internal/devices/register", postRegisterKitchenDevice(cfg))
	}

	return r
}

func postRegisterKitchenDevice(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Token") != cfg.InternalDeviceToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token"})
			return
		}
		var body struct {
			KitchenID string `json:"kitchen_id" binding:"required"`
			Label     string `json:"label"`
			CSRPEM    string `json:"csr_pem" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		kid, err := uuid.Parse(body.KitchenID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "kitchen_id"})
			return
		}
		status := "pending_signing"
		certPEM := ""
		serial := ""
		if len(cfg.DeviceCACert) > 0 && len(cfg.DeviceCAKey) > 0 {
			out, ser, err := devcert.SignCSR([]byte(body.CSRPEM), cfg.DeviceCACert, cfg.DeviceCAKey)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			certPEM = string(out)
			serial = ser
			status = "active"
		}
		id, err := cfg.Store.UpsertKitchenDevice(c.Request.Context(), kid, body.Label, body.CSRPEM, certPEM, serial, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"device_id":        id,
			"certificate_pem":  certPEM,
			"serial_number":    serial,
			"status":           status,
			"emqx_mtls_note":   "Configure EMQX listener with verify_peer and this CA; devices use issued client cert.",
		})
	}
}
