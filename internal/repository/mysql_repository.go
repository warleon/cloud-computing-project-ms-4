package repository

import (
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func (r *mysqlRepo) CreateRule(rule *rules.Rule) error {
	return r.db.Create(rule).Error
}

func (r *mysqlRepo) ReadRule(id uint) (*rules.Rule, error) {
	var rule rules.Rule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *mysqlRepo) ReadRules(size int, offset int) ([]rules.Rule, error) {
	var out []rules.Rule
	if err := r.db.Offset(offset).Limit(size).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mysqlRepo) UpdateRule(rule *rules.Rule) error {
	return r.db.Updates(rule).Error
}

func (r *mysqlRepo) DeleteRule(id uint) error {
	return r.db.Delete(&rules.Rule{}, id).Error
}

func (r *mysqlRepo) CreateAudit(audit *AuditLog) error {
	return r.db.Create(audit).Error
}

func (r *mysqlRepo) ReadAuidits(size int, offset int) ([]AuditLog, error) {
	var out []AuditLog
	if err := r.db.Offset(offset).Limit(size).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mysqlRepo) FindRulesByType(ruleType string) ([]rules.Rule, error) {
	var out []rules.Rule
	if err := r.db.Where("type = ?", ruleType).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mysqlRepo) IsAccountSanctioned(accID string) (bool, error) {
	var s rules.Sanction
	err := r.db.Where(&rules.Sanction{AccID: accID}).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &mysqlRepo{db: db}
}
