package reports

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunDailyClientAggregate_NilDeps(t *testing.T) {
	err := RunDailyClientAggregate(context.Background(), nil, nil, nil, time.Now())
	require.Error(t, err)
}
