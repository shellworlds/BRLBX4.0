package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestAudienceIncludes(t *testing.T) {
	m := jwt.MapClaims{"aud": []any{"a", "b"}}
	require.True(t, audienceIncludes(m, "b"))
	require.False(t, audienceIncludes(m, "c"))
	m2 := jwt.MapClaims{"aud": "x"}
	require.True(t, audienceIncludes(m2, "x"))
}

func TestAdminAndRolesMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin", AdminAPIKeyMiddleware("X-Admin", "secret"), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("X-Admin", "wrong")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req2.Header.Set("X-Admin", "secret")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusOK, w2.Code)

	r2 := gin.New()
	r2.GET("/need", RequireRole("admin"), func(c *gin.Context) { c.Status(200) })
	w3 := httptest.NewRecorder()
	r2.ServeHTTP(w3, httptest.NewRequest(http.MethodGet, "/need", nil))
	require.Equal(t, http.StatusForbidden, w3.Code)
}

func jwkRSAPub(t *testing.T, pub *rsa.PublicKey, kid string) map[string]any {
	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()
	return map[string]any{
		"keys": []any{map[string]any{
			"kty": "RSA",
			"kid": kid,
			"use": "sig",
			"n":   base64.RawURLEncoding.EncodeToString(nBytes),
			"e":   base64.RawURLEncoding.EncodeToString(eBytes),
		}},
	}
}

func TestValidator_ParseWithFakeJWKS(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pub := &priv.PublicKey
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/.well-known/jwks.json", r.URL.Path)
		_ = json.NewEncoder(w).Encode(jwkRSAPub(t, pub, "k1"))
	}))
	t.Cleanup(srv.Close)

	v := NewValidator(Config{Issuer: srv.URL, Audience: "api", Client: srv.Client()})
	WithIssuer(v, srv.URL)

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": srv.URL,
		"aud": "api",
		"sub": "auth0|123",
		"email": "u@example.com",
		"roles": []any{"admin"},
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tok.Header["kid"] = "k1"
	s, err := tok.SignedString(priv)
	require.NoError(t, err)

	u, err := v.Parse(s)
	require.NoError(t, err)
	require.Equal(t, "auth0|123", u.Subject)
	require.Contains(t, u.Roles, "admin")
}

func TestOptionalMiddleware_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	v := &Validator{cfg: Config{}, cache: map[string]*rsa.PublicKey{}}
	r := gin.New()
	r.Use(OptionalMiddleware(v))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestSubjectFromContext_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	require.Equal(t, "", SubjectFromContext(c))
	require.Nil(t, RolesFromContext(c))
}
