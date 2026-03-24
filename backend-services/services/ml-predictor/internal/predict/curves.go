package predict

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

type EnergyRequest struct {
	KitchenID   uuid.UUID `json:"kitchen_id"`
	HoursAhead  int       `json:"hours_ahead" binding:"required"`
	SampleKW    []float64 `json:"history_kw"` // optional stub history
}

type EnergyPoint struct {
	Hour int     `json:"hour"`
	KW   float64 `json:"kw"`
}

type EnergyResponse struct {
	KitchenID uuid.UUID     `json:"kitchen_id"`
	Points    []EnergyPoint `json:"points"`
	Note      string        `json:"note"`
}

// PredictEnergyCurve uses a lightweight trend + diurnal sinusoid as a Day-3 stand-in for ONNX/TF serving.
func PredictEnergyCurve(_ context.Context, req EnergyRequest) EnergyResponse {
	base := 12.0
	if len(req.SampleKW) > 2 {
		mean := stat.Mean(req.SampleKW, nil)
		if mean > 0 {
			base = mean
		}
	}
	slope := 0.05
	if len(req.SampleKW) > 2 {
		x := make([]float64, len(req.SampleKW))
		for i := range x {
			x[i] = float64(i)
		}
		_, beta := stat.LinearRegression(x, req.SampleKW, nil, false)
		if !math.IsNaN(beta) {
			slope = beta / float64(max(1, len(req.SampleKW)))
		}
	}
	out := make([]EnergyPoint, 0, req.HoursAhead)
	for h := 1; h <= req.HoursAhead; h++ {
		diurnal := 2 * math.Sin(2*math.Pi*float64(h)/24.0)
		kw := base + slope*float64(h) + diurnal
		if kw < 0 {
			kw = 0
		}
		out = append(out, EnergyPoint{Hour: h, KW: kw})
	}
	return EnergyResponse{KitchenID: req.KitchenID, Points: out, Note: "stub: gonum regression + diurnal; replace with trained model"}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type DemandRequest struct {
	VendorID   uuid.UUID `json:"vendor_id" binding:"required"`
	Slots      int       `json:"slots" binding:"required"` // hours
	History    []float64 `json:"history_meals"`
}

type DemandPoint struct {
	Slot      int     `json:"slot"`
	MealCount float64 `json:"meal_count"`
}

type DemandResponse struct {
	VendorID uuid.UUID     `json:"vendor_id"`
	Points   []DemandPoint `json:"points"`
	Note     string        `json:"note"`
}

// PredictDemand projects meal counts per time bucket (very coarse stub).
func PredictDemand(_ context.Context, req DemandRequest) DemandResponse {
	baseMeals := 120.0
	if len(req.History) > 0 {
		baseMeals = floats.Max(req.History) * 0.9
	}
	out := make([]DemandPoint, 0, req.Slots)
	for s := 1; s <= req.Slots; s++ {
		season := 15 * math.Sin(2*math.Pi*float64(s)/float64(max(1, req.Slots)))
		out = append(out, DemandPoint{Slot: s, MealCount: baseMeals + season})
	}
	return DemandResponse{VendorID: req.VendorID, Points: out, Note: "stub curve; ingest vendor-ecosystem history in Day-4+"}
}

type ModelInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	AccuracyMAE float64   `json:"accuracy_mae_stub"`
	Updated     time.Time `json:"updated_at"`
}

func Catalog() []ModelInfo {
	now := time.Now().UTC()
	return []ModelInfo{
		{Name: "energy_curve_stub", Version: "0.1.0", AccuracyMAE: 0.18, Updated: now},
		{Name: "vendor_demand_stub", Version: "0.1.0", AccuracyMAE: 12.5, Updated: now},
	}
}
