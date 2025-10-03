package repository

import (
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
	"gorm.io/gorm"
)

// AuditLog stores decisions for regulatory reporting.
type AuditLog struct {
	gorm.Model
	TransactionID string         `gorm:"index" json:"transactionId"`
	CustomerID    string         `gorm:"index" json:"customerId"`
	DecisionID    uint           `json:"decisionId"`
	Decision      rules.Decision `json:"decision"`
}

// Repository defines DB operations needed by the service.
type Repository interface {
	CreateRule(r *rules.Rule) error
	ReadRule(id uint) (*rules.Rule, error)
	ReadRules(size int, offset int) ([]rules.Rule, error)
	UpdateRule(r *rules.Rule) error
	DeleteRule(id uint) error
	CreateAudit(a *AuditLog) error
	ReadAuidits(size int, offset int) ([]AuditLog, error)
	// FindRulesByType returns rules filtered by their Type field (e.g. "amount_threshold").
	FindRulesByType(ruleType string) ([]rules.Rule, error)
	// IsAccountSanctioned checks whether an account identifier exists in the sanctions table.
	IsAccountSanctioned(accID string) (bool, error)
}

var RepositoryTables = []any{
	&AuditLog{},
}
