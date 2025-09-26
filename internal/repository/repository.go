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
	CreateRule(r *rules.Rule) error
	ReadRule(id uint) (*rules.Rule, error)
	ReadRules(size int, offset int) ([]rules.Rule, error)
	UpdateRule(r *rules.Rule) error
	DeleteRule(id uint) error
	CreateAudit(a *AuditLog) error
	ReadAuidits(size int, offset int) ([]AuditLog, error)
}
