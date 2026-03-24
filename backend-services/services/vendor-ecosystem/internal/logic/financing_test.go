package logic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

func TestEvaluateAdvance(t *testing.T) {
	t.Parallel()
	v := &repo.Vendor{FSSAIScore: 90}
	status, reason := EvaluateAdvance(v, 20_000)
	require.Equal(t, "approved", status)
	require.NotEmpty(t, reason)

	v2 := &repo.Vendor{FSSAIScore: 50}
	status2, _ := EvaluateAdvance(v2, 20_000)
	require.Equal(t, "rejected", status2)
}

func TestRepaymentAmount(t *testing.T) {
	t.Parallel()
	require.InDelta(t, 10.0, RepaymentAmount(100, 50), 1e-9)
	require.InDelta(t, 50.0, RepaymentAmount(1000, 50), 1e-9)
}
