package energy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClient_PostReading(t *testing.T) {
	k := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/kitchens/"+k.String()+"/readings", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)

	c := &Client{BaseURL: srv.URL, HTTP: srv.Client(), IngestToken: "t"}
	err := c.PostReading(context.Background(), k, Reading{GridPower: 1})
	require.NoError(t, err)
}

func TestClient_PostReading_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)
	c := &Client{BaseURL: srv.URL, HTTP: srv.Client()}
	err := c.PostReading(context.Background(), uuid.New(), Reading{})
	require.Error(t, err)
}
