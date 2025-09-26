package repository

import (
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
	"gorm.io/gorm"
)

// AuditLog stores decisions for regulatory reporting.
type AuditLog struct {
	gorm.Model
	TransactionID string         `gorm:"size:100;index" json:"transactionId"`
	CustomerID    string         `gorm:"size:100;index" json:"customerId"`
	Decision      rules.Decision `gorm:"size:50" json:"decision"`
}

// Repository defines DB operations needed by the service.
type Repository interface {
	CreateRule(r *rules.RuleListItem) error
	ReadRule(id uint) (rules.RuleListItem, error)
	ReadRules() ([]rules.RuleListItem, error)
	UpdateRule(r rules.RuleListItem) error
	DeleteRule(id uint) error
}
