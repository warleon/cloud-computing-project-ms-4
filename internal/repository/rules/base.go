package rules

import (
	"time"

	"github.com/warleon/ms4-compliance-service/internal/dto"
	"gorm.io/gorm"
)

type RuleType string

const (
	RuleTypeAmountThreshold RuleType = "amount_threshold"
	RuleTypeSanctionsList   RuleType = "sanctions_list"
)

// Rule represents a compliance rule stored in DB.
type RuleBase struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        RuleType  `gorm:"type:enum('amount_threshold','sanctions_list');not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ComplianceRule is the interface each rule implements.
type ComplianceRule interface {
	Validate(tx dto.Transaction) Decision
}
type RuleExtras struct {
	Threshold *float64
	DB        *gorm.DB
}

var RuleTables = []any{
	&AmountThresholdRule{},
	&BlacklistRule{},
	&Sanction{},
}
