package repository

import (
	"time"

	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
)

// AuditLog stores decisions for regulatory reporting.
type AuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID string    `gorm:"size:100;index" json:"transactionId"`
	CustomerID    string    `gorm:"size:100;index" json:"customerId"`
	Decision      string    `gorm:"size:50" json:"decision"` // approve, review, reject
	Reason        string    `gorm:"type:text" json:"reason"`
	Meta          string    `gorm:"type:text" json:"meta"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Repository defines DB operations needed by the service.
type Repository interface {
	CreateRule(r *rules.RuleBase, extras rules.RuleExtras) error
	ListRules() ([]rules.RuleBase, error)
	FindRulesByType(t string) ([]rules.RuleBase, error)
	CreateAudit(a *AuditLog) error
	ListAuditsForCustomer(customerID string, limit int) ([]AuditLog, error)
}
