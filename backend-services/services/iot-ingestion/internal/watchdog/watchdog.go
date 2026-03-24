package watchdog

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Tracker records last telemetry time per kitchen (in-memory; extend with Redis for HA).
type Tracker struct {
	mu      sync.RWMutex
	last    map[uuid.UUID]time.Time
	offline map[uuid.UUID]time.Time
}

func NewTracker() *Tracker {
	return &Tracker{
		last:    map[uuid.UUID]time.Time{},
		offline: map[uuid.UUID]time.Time{},
	}
}

func (t *Tracker) Seen(kitchen uuid.UUID, ts time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last[kitchen] = ts
	delete(t.offline, kitchen)
}

func (t *Tracker) Check(ctx context.Context, stale time.Duration, onOffline func(ctx context.Context, kitchen uuid.UUID, last time.Time) error) error {
	now := time.Now().UTC()
	type candidate struct {
		k    uuid.UUID
		last time.Time
	}

	t.mu.Lock()
	var pending []candidate
	for kitchen, last := range t.last {
		if now.Sub(last) <= stale {
			continue
		}
		if _, already := t.offline[kitchen]; already {
			continue
		}
		t.offline[kitchen] = now
		pending = append(pending, candidate{k: kitchen, last: last})
	}
	t.mu.Unlock()

	for _, c := range pending {
		if err := onOffline(ctx, c.k, c.last); err != nil {
			return err
		}
	}
	return nil
}
