package repository

import (
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) CreateRule(rule *Rule) error {
	return r.db.Create(rule).Error
}

func (r *mysqlRepo) ListRules() ([]Rule, error) {
	var out []Rule
	if err := r.db.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mysqlRepo) FindRulesByType(t string) ([]Rule, error) {
	var out []Rule
	if err := r.db.Where("type = ?", t).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mysqlRepo) CreateAudit(a *AuditLog) error {
	return r.db.Create(a).Error
}

func (r *mysqlRepo) ListAuditsForCustomer(customerID string, limit int) ([]AuditLog, error) {
	var out []AuditLog
	if err := r.db.Where("customer_id = ?", customerID).Order("created_at desc").Limit(limit).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
