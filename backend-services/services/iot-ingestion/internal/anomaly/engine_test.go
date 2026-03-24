package anomaly

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEngineObserveNoFireEarly(t *testing.T) {
	e := New(5, 3)
	id := uuid.New()
	fire, _ := e.Observe(id, 10)
	require.False(t, fire)
}

func TestEngineObserveSpike(t *testing.T) {
	e := New(12, 2)
	id := uuid.New()
	for i := 0; i < 8; i++ {
		e.Observe(id, 10+float64(i)*0.01) // tight cluster
	}
	fire, why := e.Observe(id, 500)
	require.True(t, fire)
	require.NotEmpty(t, why)
}
