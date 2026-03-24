package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/logic"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

func NewRouter(store *repo.Store) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", metrics.Handler())

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

	return r
}
