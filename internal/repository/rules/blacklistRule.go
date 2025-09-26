package rules

import (
	"github.com/warleon/ms4-compliance-service/internal/dto"
	"gorm.io/gorm"
)

type Sanction struct {
	gorm.Model
	AccID string `gorm:"size:100;index"` // account/customer identifier
}

type BlacklistRule struct {
	RuleBase
	db *gorm.DB
}

func (r *BlacklistRule) Validate(tx dto.Transaction) Decision {
	var entry Sanction

	// Fast lookup using indexed column
	err := r.db.Where(&Sanction{AccID: tx.FromAcc}).Or(&Sanction{AccID: tx.ToAcc}).First(&entry).Error

	if err == nil {
		return Decision{
			Approved: false,
			Reason:   "Account is blacklisted",
		}
	}

	return Decision{Approved: true, Reason: "OK"}
}
