package reports

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunDailyVendorAggregate_Nil(t *testing.T) {
	err := RunDailyVendorAggregate(context.Background(), nil, nil, time.Now())
	require.Error(t, err)
}
