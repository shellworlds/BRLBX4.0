package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ctxKeySubject = "auth_subject"
const ctxKeyRoles = "auth_roles"
const ctxKeyEmail = "auth_email"

// User is the normalized identity extracted from a verified JWT.
type User struct {
	Subject string
	Email   string
	Roles   []string
}

// Config holds JWKS discovery and audience issuer checks.
type Config struct {
	Domain   string
	Audience string
	Issuer   string
	Client   *http.Client
}

// JWKS minimal model.
type jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

// Validator verifies RS256 JWTs against Auth0 JWKS.
type Validator struct {
	cfg     Config
	mu      sync.RWMutex
	jwksURL string
	cache   map[string]*rsa.PublicKey
}

// NewValidator builds a validator for an Auth0 tenant domain (e.g. dev-xyz.us.auth0.com).
func NewValidator(cfg Config) *Validator {
	if cfg.Client == nil {
		cfg.Client = http.DefaultClient
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "https://" + strings.TrimPrefix(cfg.Domain, "https://")
	}
	return &Validator{
		cfg:     cfg,
		jwksURL: strings.TrimSuffix(cfg.Issuer, "/") + "/.well-known/jwks.json",
		cache:   map[string]*rsa.PublicKey{},
	}
}

func (v *Validator) keyFunc(token *jwt.Token) (any, error) {
	if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
		return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
	}
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return nil, errors.New("missing kid")
	}
	key, err := v.getKey(context.Background(), kid)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (v *Validator) getKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	if k, ok := v.cache[kid]; ok {
		v.mu.RUnlock()
		return k, nil
	}
	v.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := v.cfg.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	var doc jwks
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, j := range doc.Keys {
		if j.Kty != "RSA" || j.N == "" || j.E == "" {
			continue
		}
		pub, err := rsaKeyFromComponents(j.N, j.E)
		if err != nil {
			continue
		}
		v.cache[j.Kid] = pub
	}
	if k, ok := v.cache[kid]; ok {
		return k, nil
	}
	return nil, errors.New("jwks: key not found")
}

func rsaKeyFromComponents(nb64, eb64 string) (*rsa.PublicKey, error) {
	decode := func(s string) ([]byte, error) {
		return base64.RawURLEncoding.DecodeString(s)
	}
	nb, err := decode(nb64)
	if err != nil {
		return nil, err
	}
	eb, err := decode(eb64)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nb)
	e := int(new(big.Int).SetBytes(eb).Int64())
	if e == 0 {
		e = 65537
	}
	return &rsa.PublicKey{N: n, E: e}, nil
}

// Parse validates token string and returns a normalized user.
func (v *Validator) Parse(tokenStr string) (*User, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	tok, err := parser.Parse(tokenStr, v.keyFunc)
	if err != nil || tok == nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	m, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if v.cfg.Issuer != "" {
		iss, _ := m["iss"].(string)
		if iss != v.cfg.Issuer {
			return nil, errors.New("invalid issuer")
		}
	}
	if v.cfg.Audience != "" && !audienceIncludes(m, v.cfg.Audience) {
		return nil, errors.New("invalid audience")
	}
	sub, _ := m["sub"].(string)
	email, _ := m["email"].(string)
	roles := rolesFromMap(m)
	return &User{Subject: sub, Email: email, Roles: roles}, nil
}

func audienceIncludes(m jwt.MapClaims, want string) bool {
	raw, ok := m["aud"]
	if !ok {
		return false
	}
	switch t := raw.(type) {
	case string:
		return t == want
	case []any:
		for _, v := range t {
			if s, ok := v.(string); ok && s == want {
				return true
			}
		}
	}
	return false
}

func rolesFromMap(m jwt.MapClaims) []string {
	out := map[string]struct{}{}
	if raw, ok := m["roles"]; ok {
		mergeRolesJSON(out, raw)
	}
	for k, v := range m {
		if strings.HasSuffix(k, "/roles") {
			mergeRolesJSON(out, v)
		}
	}
	var list []string
	for r := range out {
		list = append(list, r)
	}
	return list
}

// Middleware returns gin middleware enforcing Bearer JWT.
func Middleware(v *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(h[7:])
		user, err := v.Parse(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ctxKeySubject, user.Subject)
		c.Set(ctxKeyEmail, user.Email)
		c.Set(ctxKeyRoles, user.Roles)
		c.Next()
	}
}

// OptionalMiddleware parses JWT if present; continues as anonymous otherwise.
func OptionalMiddleware(v *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.Next()
			return
		}
		raw := strings.TrimSpace(h[7:])
		user, err := v.Parse(raw)
		if err != nil {
			c.Next()
			return
		}
		c.Set(ctxKeySubject, user.Subject)
		c.Set(ctxKeyEmail, user.Email)
		c.Set(ctxKeyRoles, user.Roles)
		c.Next()
	}
}

// RequireRole returns middleware that requires one of roles when JWT was validated.
func RequireRole(roles ...string) gin.HandlerFunc {
	want := map[string]struct{}{}
	for _, r := range roles {
		want[strings.ToLower(strings.TrimSpace(r))] = struct{}{}
	}
	return func(c *gin.Context) {
		raw, ok := c.Get(ctxKeyRoles)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing roles"})
			return
		}
		list, _ := raw.([]string)
		for _, r := range list {
			if _, ok := want[strings.ToLower(strings.TrimSpace(r))]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}

// SubjectFromContext returns Auth0 subject or empty.
func SubjectFromContext(c *gin.Context) string {
	if v, ok := c.Get(ctxKeySubject); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RolesFromContext returns roles slice (may be nil).
func RolesFromContext(c *gin.Context) []string {
	if v, ok := c.Get(ctxKeyRoles); ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return nil
}

// AdminAPIKeyMiddleware allows admin routes via static key (bootstrap / machine admin).
func AdminAPIKeyMiddleware(headerName, expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if expected == "" {
			c.Next()
			return
		}
		if c.GetHeader(headerName) != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "admin key required"})
			return
		}
		c.Next()
	}
}

func mergeRolesJSON(dst map[string]struct{}, raw any) {
	switch t := raw.(type) {
	case []any:
		for _, v := range t {
			if s, ok := v.(string); ok {
				dst[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
			}
		}
	case string:
		dst[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}
}

// WithIssuer overrides issuer for tests.
func WithIssuer(v *Validator, issuer string) {
	v.cfg.Issuer = issuer
	parsed, err := url.Parse(issuer)
	if err == nil && parsed != nil && parsed.Scheme != "" {
		v.jwksURL = strings.TrimSuffix(issuer, "/") + "/.well-known/jwks.json"
		return
	}
	v.jwksURL = strings.TrimSuffix(issuer, "/") + "/.well-known/jwks.json"
}
