package rules

import (
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
	gorm.Model
	Name        string   `gorm:"size:255;not null" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	Type        RuleType `gorm:"type:enum('amount_threshold','sanctions_list');not null"`
	User        string   `gorm:"index"`
}

// ComplianceRule is the interface each rule implements.
type ComplianceRule interface {
	Validate(tx dto.Transaction) Decision
}
type RuleExtras struct {
	Threshold *float64
	DB        *gorm.DB
}

type Rule struct {
	RuleBase
	RuleExtras
}

var RuleTables = []any{
	&Rule{},
	&Sanction{},
}
