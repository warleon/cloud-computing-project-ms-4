package service

import (
	"context"

	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
)

// Rule CRUD helpers on the service layer. These simply delegate to the repository
// and exist to provide a clear service boundary and allow future business logic
// to be applied around rule operations.

// CreateRule inserts a new rule record.
func (s *ComplianceService) CreateRule(ctx context.Context, r *rules.Rule) error {
	return s.Repo.CreateRule(r)
}

// GetRule returns a single rule by ID.
func (s *ComplianceService) GetRule(ctx context.Context, id uint) (*rules.Rule, error) {
	return s.Repo.ReadRule(id)
}

// ListRules returns a paginated list of rules.
func (s *ComplianceService) ListRules(ctx context.Context, size int, offset int) ([]rules.Rule, error) {
	return s.Repo.ReadRules(size, offset)
}

// UpdateRule updates an existing rule.
func (s *ComplianceService) UpdateRule(ctx context.Context, r *rules.Rule) error {
	return s.Repo.UpdateRule(r)
}

// DeleteRule removes a rule by ID.
func (s *ComplianceService) DeleteRule(ctx context.Context, id uint) error {
	return s.Repo.DeleteRule(id)
}
