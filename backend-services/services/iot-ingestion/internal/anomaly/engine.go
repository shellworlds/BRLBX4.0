package anomaly

import (
	"fmt"
	"math"
	"sync"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat"
)

// Engine maintains a short sliding window of grid power (kW) per kitchen and flags 3-sigma excursions.
type Engine struct {
	mu       sync.Mutex
	window   int
	sigma    float64
	history  map[uuid.UUID][]float64
	minSamples int
}

func New(window int, sigma float64) *Engine {
	if window < 5 {
		window = 5
	}
	if sigma <= 0 {
		sigma = 3
	}
	return &Engine{
		window:     window,
		sigma:      sigma,
		history:    map[uuid.UUID][]float64{},
		minSamples: 5,
	}
}

// Observe ingests a new sample; returns true if the latest reading is anomalous vs prior window.
func (e *Engine) Observe(kitchen uuid.UUID, gridKW float64) (bool, string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	h := append(e.history[kitchen], gridKW)
	if len(h) > e.window {
		h = h[len(h)-e.window:]
	}
	e.history[kitchen] = h

	if len(h) < e.minSamples {
		return false, ""
	}
	prior := h[:len(h)-1]
	current := h[len(h)-1]
	mean := stat.Mean(prior, nil)
	stdev := stat.StdDev(prior, nil)
	if stdev == 0 {
		return false, ""
	}
	z := math.Abs(current-mean) / stdev
	if z >= e.sigma {
		return true, fmt.Sprintf("grid_power z-score %.2f (>%.2f sigma) mean=%.3f stdev=%.3f current=%.3f", z, e.sigma, mean, stdev, current)
	}
	return false, ""
}
