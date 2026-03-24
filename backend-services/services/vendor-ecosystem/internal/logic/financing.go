package logic

import (
	"math"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

// EvaluateAdvance is an ML stub: uses average transaction volume proxy and FSSAI score.
func EvaluateAdvance(v *repo.Vendor, avgVolume float64) (status string, reason string) {
	const minVol = 10_000
	const minScore = 80
	if v.FSSAIScore < minScore {
		return "rejected", "fssai_score_below_threshold"
	}
	if avgVolume < minVol {
		return "rejected", "avg_volume_below_threshold"
	}
	return "approved", "rules_engine_stub_pass"
}

// RepaymentAmount computes automatic repayment for a transaction (10% capped by remaining balance).
func RepaymentAmount(transactionAmt, remaining float64) float64 {
	repayment := transactionAmt * 0.10
	if repayment > remaining {
		return remaining
	}
	if repayment < 0 {
		return 0
	}
	return repayment
}

// ApplyRepayment returns new remaining balance after subtracting repayment.
func ApplyRepayment(remaining, repayment float64) float64 {
	n := remaining - repayment
	if n < 0 {
		return 0
	}
	// avoid -0 oddities
	return math.Round(n*100) / 100
}
