package rules

import "github.com/warleon/ms4-compliance-service/internal/dto"

type AmountThresholdRule struct {
	RuleBase
	Threshold float64
}

func (r *AmountThresholdRule) Validate(tx dto.Transaction) Decision {
	if tx.Amount > r.Threshold {
		return Decision{
			Approved: false,
			Reason:   "Transaction exceeds threshold",
		}
	}
	return Decision{Approved: true, Reason: "OK"}
}
