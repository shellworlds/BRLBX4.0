package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/paymentintent"
	stripeSub "github.com/stripe/stripe-go/v81/subscription"
	"github.com/stripe/stripe-go/v81/webhook"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/payments/internal/repo"
)

type RouterConfig struct {
	Store           *repo.Store
	Validator       *auth.Validator
	StripeSecretKey string
	StripeWebhook   string
	DefaultPriceID  string
	InternalToken   string
	EnableSwagger   bool
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	cors.Use(r)
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", metrics.Handler())

	if cfg.StripeSecretKey != "" {
		stripe.Key = cfg.StripeSecretKey
	}

	v1 := r.Group("/api/v1")

	v1.POST("/webhooks/stripe", func(c *gin.Context) {
		if cfg.StripeWebhook == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "webhook secret not configured"})
			return
		}
		payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "body"})
			return
		}
		sig := c.GetHeader("Stripe-Signature")
		ev, err := webhook.ConstructEvent(payload, sig, cfg.StripeWebhook)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "signature"})
			return
		}
		dup, err := cfg.Store.TryMarkWebhookEvent(c.Request.Context(), ev.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if dup {
			c.JSON(http.StatusOK, gin.H{"status": "duplicate"})
			return
		}
		switch ev.Type {
		case "customer.subscription.updated", "customer.subscription.deleted":
			var sub stripe.Subscription
			if err := json.Unmarshal(ev.Data.Raw, &sub); err == nil && sub.ID != "" {
				status := string(sub.Status)
				if ev.Type == "customer.subscription.deleted" {
					status = "canceled"
				}
				var next *time.Time
				if sub.CurrentPeriodEnd > 0 {
					t := time.Unix(sub.CurrentPeriodEnd, 0).UTC()
					next = &t
				}
				_ = cfg.Store.UpdateSubscriptionByStripeID(c.Request.Context(), sub.ID, status, next)
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if cfg.Validator != nil {
		bill := v1.Group("")
		bill.Use(auth.Middleware(cfg.Validator))
		bill.Use(auth.RequireRole("client", "admin"))

		bill.POST("/subscriptions", postSubscription(cfg))
		bill.GET("/subscriptions/me", getSubscriptionMe(cfg))
		bill.POST("/subscriptions/me/cancel", cancelSubscription(cfg))
		bill.POST("/carbon-credits/buy", postCarbon(cfg))
	}

	v1.POST("/transactions/record", postInternalMeal(cfg))

	return r
}

func postSubscription(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := auth.StringClaimFromContext(c, "https://borelsigma.com/client_id", "client_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id claim required"})
			return
		}
		var body struct {
			Plan string `json:"plan" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sub := &repo.Subscription{ClientID: cid, Plan: body.Plan, Status: "stub_active"}
		if cfg.StripeSecretKey != "" && cfg.DefaultPriceID != "" {
			email := auth.StringClaimFromContext(c, "email")
			cust, err := customer.New(&stripe.CustomerParams{Email: stripe.String(email)})
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			sid := cust.ID
			st, err := stripeSub.New(&stripe.SubscriptionParams{
				Customer: stripe.String(sid),
				Items: []*stripe.SubscriptionItemsParams{
					{Price: stripe.String(cfg.DefaultPriceID)},
				},
			})
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			sub.Status = string(st.Status)
			sub.StripeCustomerID = &sid
			ss := st.ID
			sub.StripeSubscriptionID = &ss
			if st.CurrentPeriodEnd > 0 {
				t := time.Unix(st.CurrentPeriodEnd, 0).UTC()
				sub.NextBilling = &t
			}
		}
		if err := cfg.Store.InsertSubscription(c.Request.Context(), sub); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, sub)
	}
}

func getSubscriptionMe(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := auth.StringClaimFromContext(c, "https://borelsigma.com/client_id", "client_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id claim required"})
			return
		}
		s, err := cfg.Store.GetActiveSubscription(c.Request.Context(), cid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if s == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no subscription"})
			return
		}
		c.JSON(http.StatusOK, s)
	}
}

func cancelSubscription(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := auth.StringClaimFromContext(c, "https://borelsigma.com/client_id", "client_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id claim required"})
			return
		}
		if cfg.StripeSecretKey != "" {
			s, err := cfg.Store.GetActiveSubscription(c.Request.Context(), cid)
			if err == nil && s != nil && s.StripeSubscriptionID != nil && *s.StripeSubscriptionID != "" {
				_, _ = stripeSub.Cancel(*s.StripeSubscriptionID, nil)
			}
		}
		n, err := cfg.Store.CancelSubscription(c.Request.Context(), cid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"canceled": n})
	}
}

func postCarbon(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := auth.StringClaimFromContext(c, "https://borelsigma.com/client_id", "client_id")
		if cid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id claim required"})
			return
		}
		var body struct {
			Tonnes float64 `json:"tonnes" binding:"required"`
			Amount float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rec := &repo.CarbonPurchase{ClientID: cid, Tonnes: body.Tonnes, Amount: body.Amount}
		if cfg.StripeSecretKey != "" {
			amountCents := int64(body.Amount * 100)
			if amountCents < 50 {
				amountCents = 50
			}
			pi, err := paymentintent.New(&stripe.PaymentIntentParams{
				Amount:   stripe.Int64(amountCents),
				Currency: stripe.String("usd"),
				AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
					Enabled: stripe.Bool(true),
				},
			})
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			s := pi.ID
			rec.StripePaymentIntentID = &s
		}
		if err := cfg.Store.InsertCarbon(c.Request.Context(), rec); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, rec)
	}
}

func postInternalMeal(cfg RouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.InternalToken == "" || c.GetHeader("X-Internal-Token") != cfg.InternalToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token"})
			return
		}
		var body struct {
			VendorID              uuid.UUID `json:"vendor_id" binding:"required"`
			KitchenID             uuid.UUID `json:"kitchen_id" binding:"required"`
			MealCount             int       `json:"meal_count" binding:"required"`
			Amount                float64   `json:"amount" binding:"required"`
			PaymentMethod         string    `json:"payment_method"`
			StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rec := &repo.MealPaymentRecord{
			VendorID:      body.VendorID,
			KitchenID:     body.KitchenID,
			MealCount:     body.MealCount,
			Amount:        body.Amount,
			PaymentMethod: nonemptyPtr(body.PaymentMethod),
		}
		if body.StripePaymentIntentID != "" {
			rec.StripePaymentIntentID = &body.StripePaymentIntentID
		}
		if err := cfg.Store.InsertMealRecord(c.Request.Context(), rec); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, rec)
	}
}

func nonemptyPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
