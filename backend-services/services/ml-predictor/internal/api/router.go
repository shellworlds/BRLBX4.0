package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/internal/cache"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/ml-predictor/internal/predict"
)

var predictLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "ml_predict_latency_seconds",
	Help:    "Prediction handler latency",
	Buckets: prometheus.DefBuckets,
}, []string{"route"})

type RouterConfig struct {
	Cache         cache.PredictCache
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

	if cfg.Cache == nil {
		cfg.Cache = cache.Noop{}
	}

	v1 := r.Group("/api/v1")

	v1.GET("/predict/models", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": predict.Catalog()})
	})

	v1.POST("/predict/energy", func(c *gin.Context) {
		start := time.Now()
		defer func() {
			predictLatency.WithLabelValues("energy").Observe(time.Since(start).Seconds())
		}()
		var body predict.EnergyRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if body.HoursAhead <= 0 || body.HoursAhead > 168 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hours_ahead"})
			return
		}
		key := "energy:" + body.KitchenID.String() + ":" + strconv.Itoa(body.HoursAhead)
		var cached predict.EnergyResponse
		if cfg.Cache.Get(c.Request.Context(), key, &cached) {
			c.JSON(http.StatusOK, cached)
			return
		}
		out := predict.PredictEnergyCurve(c.Request.Context(), body)
		_ = cfg.Cache.Set(c.Request.Context(), key, out, 5*time.Minute)
		c.JSON(http.StatusOK, out)
	})

	v1.POST("/predict/demand", func(c *gin.Context) {
		start := time.Now()
		defer func() {
			predictLatency.WithLabelValues("demand").Observe(time.Since(start).Seconds())
		}()
		var body predict.DemandRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if body.Slots <= 0 || body.Slots > 168 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "slots"})
			return
		}
		out := predict.PredictDemand(c.Request.Context(), body)
		c.JSON(http.StatusOK, out)
	})

	return r
}
