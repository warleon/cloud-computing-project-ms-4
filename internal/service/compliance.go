package service

import (
	"context"

	"github.com/warleon/ms4-compliance-service/internal/dto"
	"github.com/warleon/ms4-compliance-service/internal/repository"
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
)

// ComplianceService contains business logic.
type ComplianceService struct {
	Repo repository.Repository
}

func NewComplianceService(repo repository.Repository) *ComplianceService {
	return &ComplianceService{Repo: repo}
}

func (s *ComplianceService) ValidateTransaction(ctx context.Context, in dto.Transaction) (*rules.Decision, error) {
	// 1) Evaluate amount threshold rules
	amtRules, err := s.Repo.FindRulesByType(string(rules.RuleTypeAmountThreshold))
	if err != nil {
		return nil, err
	}

	for _, r := range amtRules {
		if r.Threshold == nil {
			// skip malformed rule
			continue
		}
		// create concrete rule and validate
		ar := rules.AmountThresholdRule{
			RuleBase:  r.RuleBase,
			Threshold: *r.Threshold,
		}
		dec := ar.Validate(in)
		if !dec.Approved {
			return &dec, nil
		}
	}

	// 2) Evaluate sanctions / blacklist rules: check if either account is sanctioned
	// We rely on repository-level helper to check sanctions table quickly.
	fromSanctioned, err := s.Repo.IsAccountSanctioned(in.FromAcc)
	if err != nil {
		return nil, err
	}
	if fromSanctioned {
		d := rules.Decision{Approved: false, Reason: "From account is sanctioned"}
		return &d, nil
	}

	toSanctioned, err := s.Repo.IsAccountSanctioned(in.ToAcc)
	if err != nil {
		return nil, err
	}
	if toSanctioned {
		d := rules.Decision{Approved: false, Reason: "To account is sanctioned"}
		return &d, nil
	}

	// All checks passed
	ok := rules.Decision{Approved: true, Reason: "OK"}
	return &ok, nil

}
