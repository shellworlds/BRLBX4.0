package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/logic"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/reports"
)

type RouterConfig struct {
	Store         *repo.Store
	Daily         *repo.DailyVendorStore
	InternalToken string
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

	store := cfg.Store
	v1 := r.Group("/api/v1")

	v1.POST("/vendors", func(c *gin.Context) {
		var body struct {
			Name       string `json:"name" binding:"required"`
			FSSAIScore int    `json:"fssai_score" binding:"required"`
			Location   string `json:"location" binding:"required"`
			Contact    string `json:"contact" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := &repo.Vendor{Name: body.Name, FSSAIScore: body.FSSAIScore, Location: body.Location, Contact: body.Contact}
		if err := store.CreateVendor(c.Request.Context(), v); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	v1.GET("/vendors/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vendor id"})
			return
		}
		v, err := store.GetVendor(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	v1.POST("/vendors/:id/financing", func(c *gin.Context) {
		vid, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vendor id"})
			return
		}
		var body struct {
			Amount float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v, err := store.GetVendor(c.Request.Context(), vid)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "vendor not found"})
			return
		}
		to := time.Now().UTC()
		from := to.AddDate(0, -3, 0)
		avg, err := store.AvgTransactionVolume(c.Request.Context(), vid, from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		status, reason := logic.EvaluateAdvance(v, avg)
		f := &repo.Financing{
			VendorID:          vid,
			Amount:            body.Amount,
			Status:            status,
			RepaymentSchedule: "auto_pct_10",
		}
		if status == "approved" {
			f.RemainingBalance = body.Amount
		}
		if err := store.CreateFinancing(c.Request.Context(), f); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"financing": f, "reason": reason})
	})

	v1.GET("/vendors/:id/financing", func(c *gin.Context) {
		vid, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vendor id"})
			return
		}
		items, err := store.ListFinancing(c.Request.Context(), vid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	v1.POST("/transactions", func(c *gin.Context) {
		var body struct {
			VendorID  uuid.UUID `json:"vendor_id" binding:"required"`
			KitchenID uuid.UUID `json:"kitchen_id" binding:"required"`
			Amount    float64   `json:"amount" binding:"required"`
			MealCount int       `json:"meal_count"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := &repo.Transaction{VendorID: body.VendorID, KitchenID: body.KitchenID, Amount: body.Amount, MealCount: body.MealCount}
		if err := store.InsertTransaction(c.Request.Context(), tx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fin, err := store.LatestOpenFinancing(c.Request.Context(), body.VendorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if fin != nil {
			repay := logic.RepaymentAmount(body.Amount, fin.RemainingBalance)
			newRem := logic.ApplyRepayment(fin.RemainingBalance, repay)
			if err := store.UpdateFinancingBalance(c.Request.Context(), fin.ID, newRem); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{
				"transaction":         tx,
				"repayment_applied":   repay,
				"financing_remaining": newRem,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"transaction": tx})
	})

	v1.GET("/vendors/:id/transactions", func(c *gin.Context) {
		vid, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vendor id"})
			return
		}
		limit := 100
		if s := c.Query("limit"); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				limit = v
			}
		}
		items, err := store.ListTransactions(c.Request.Context(), vid, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	if cfg.Daily != nil {
		v1.GET("/reports/vendor/:vendor_id", func(c *gin.Context) {
			vid, err := reports.ParseVendorID(c.Param("vendor_id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "vendor_id"})
				return
			}
			fromS, toS := c.Query("from"), c.Query("to")
			if fromS == "" || toS == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "from and to required (YYYY-MM-DD)"})
				return
			}
			from, err := time.Parse("2006-01-02", fromS)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "from"})
				return
			}
			to, err := time.Parse("2006-01-02", toS)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "to"})
				return
			}
			items, err := cfg.Daily.List(c.Request.Context(), vid, from, to)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": items, "dashboard_hints": gin.H{"charts": []string{"meals", "revenue", "efficiency_stub"}}})
		})

		v1.POST("/internal/aggregate/daily", func(c *gin.Context) {
			if cfg.InternalToken == "" {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal aggregate disabled"})
				return
			}
			if c.GetHeader("X-Internal-Token") != cfg.InternalToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token"})
				return
			}
			day := time.Now().UTC().AddDate(0, 0, -1)
			if d := c.Query("day"); d != "" {
				parsed, err := time.Parse("2006-01-02", d)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "day"})
					return
				}
				day = parsed
			}
			if err := reports.RunDailyVendorAggregate(c.Request.Context(), store, cfg.Daily, day); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok", "day": day.Format("2006-01-02")})
		})
	}

	return r
}
