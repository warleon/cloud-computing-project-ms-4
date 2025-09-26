package rules

import "gorm.io/gorm"

type Decision struct {
	gorm.Model
	Approved bool
	Reason   string
}
