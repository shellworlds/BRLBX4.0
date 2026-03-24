package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/auth"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/cors"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/mail"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/metrics"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/compliance/internal/repo"
)

type RouterConfig struct {
	Store           *repo.Store
	Validator       *auth.Validator
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPass        string
	MailFrom        string
	SalesRecipients string
	EnableSwagger   bool
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	cors.Use(r)
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", metrics.Handler())

	v1 := r.Group("/api/v1")

	v1.POST("/contact", func(c *gin.Context) {
		var body struct {
			Name    string `json:"name" binding:"required"`
			Email   string `json:"email" binding:"required"`
			Company string `json:"company"`
			Message string `json:"message" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := cfg.Store.InsertContact(c.Request.Context(), body.Name, body.Email, body.Company, body.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_ = cfg.Store.InsertAudit(c.Request.Context(), "", "contact_form", "contact", id, nil)
		if cfg.SMTPHost != "" && cfg.MailFrom != "" && cfg.SalesRecipients != "" {
			subj := "Borel Sigma contact: " + body.Name
			msg := "From: " + body.Email + "\n" + body.Message
			_ = mail.SendSMTP(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.MailFrom,
				[]string{cfg.SalesRecipients}, subj, msg)
			_ = mail.SendSMTP(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.MailFrom,
				[]string{body.Email}, "We received your message", "Thank you for contacting Borel Sigma.")
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "id": id})
	})

	if cfg.Validator != nil {
		authg := v1.Group("")
		authg.Use(auth.Middleware(cfg.Validator))

		authg.POST("/consent", func(c *gin.Context) {
			var body struct {
				PolicyVersion string `json:"policy_version" binding:"required"`
				Accepted      bool   `json:"accepted" binding:"required"`
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			sub := auth.SubjectFromContext(c)
			if err := cfg.Store.InsertConsent(c.Request.Context(), sub, body.PolicyVersion, body.Accepted); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			meta, _ := json.Marshal(gin.H{"policy_version": body.PolicyVersion})
			_ = cfg.Store.InsertAudit(c.Request.Context(), sub, "consent_recorded", "user", sub, meta)
			c.JSON(http.StatusCreated, gin.H{"status": "ok"})
		})

		authg.GET("/consent/status", func(c *gin.Context) {
			sub := auth.SubjectFromContext(c)
			ver, acc, ok, err := cfg.Store.LatestConsent(c.Request.Context(), sub)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !ok {
				c.JSON(http.StatusOK, gin.H{"has_consent": false})
				return
			}
			c.JSON(http.StatusOK, gin.H{"has_consent": true, "policy_version": ver, "accepted": acc})
		})

		authg.DELETE("/users/me/data", func(c *gin.Context) {
			sub := auth.SubjectFromContext(c)
			if err := cfg.Store.RequestDeletion(c.Request.Context(), sub); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			meta, _ := json.Marshal(gin.H{"subject": sub})
			_ = cfg.Store.InsertAudit(c.Request.Context(), sub, "gdpr_deletion_requested", "user", sub, meta)
			c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
		})

		adm := authg.Group("")
		adm.Use(auth.RequireRole("admin"))
		adm.GET("/audit", func(c *gin.Context) {
			items, err := cfg.Store.ListAudit(c.Request.Context(), 200)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": items})
		})
	}

	return r
}
