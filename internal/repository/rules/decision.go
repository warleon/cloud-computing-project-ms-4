package rules

import "gorm.io/gorm"

type Decision struct {
	gorm.Model `swaggerignore:"true"`
	Approved   bool
	Reason     string
}
