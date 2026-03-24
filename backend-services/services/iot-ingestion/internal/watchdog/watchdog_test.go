package watchdog

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTracker_CheckOffline(t *testing.T) {
	tr := NewTracker()
	id := uuid.New()
	tr.Seen(id, time.Now().UTC().Add(-10*time.Minute))
	var called bool
	err := tr.Check(context.Background(), time.Minute, func(ctx context.Context, kitchen uuid.UUID, last time.Time) error {
		called = true
		require.Equal(t, id, kitchen)
		return nil
	})
	require.NoError(t, err)
	require.True(t, called)
}
