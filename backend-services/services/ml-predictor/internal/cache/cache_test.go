package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNoop(t *testing.T) {
	var n Noop
	require.False(t, n.Get(context.Background(), "k", &struct{}{}))
	require.NoError(t, n.Set(context.Background(), "k", 1, time.Minute))
}

func TestRedisDisabled(t *testing.T) {
	r := NewRedis("", "", 0, "")
	require.False(t, r.Get(context.Background(), "k", &struct{}{}))
	require.NoError(t, r.Set(context.Background(), "k", 1, time.Minute))
}
