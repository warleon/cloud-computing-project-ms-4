package repository

import (
	"time"
)

// Rule represents a compliance rule stored in DB.
type Rule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"size:100" json:"type"`    // e.g., "amount_threshold", "sanctions_list"
	Params      string    `gorm:"type:text" json:"params"` // JSON parameters
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

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
	CreateRule(r *Rule) error
	ListRules() ([]Rule, error)
	FindRulesByType(t string) ([]Rule, error)
	CreateAudit(a *AuditLog) error
	ListAuditsForCustomer(customerID string, limit int) ([]AuditLog, error)
}
