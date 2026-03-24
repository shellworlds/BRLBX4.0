package predict

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPredictEnergyCurve(t *testing.T) {
	out := PredictEnergyCurve(context.Background(), EnergyRequest{
		KitchenID:  uuid.New(),
		HoursAhead: 3,
		SampleKW:   []float64{10, 11, 12},
	})
	require.Len(t, out.Points, 3)
	require.Greater(t, out.Points[0].KW, 0.0)
}

func TestPredictDemand(t *testing.T) {
	out := PredictDemand(context.Background(), DemandRequest{VendorID: uuid.New(), Slots: 2, History: []float64{100, 110}})
	require.Len(t, out.Points, 2)
}

func TestCatalog(t *testing.T) {
	require.NotEmpty(t, Catalog())
}
