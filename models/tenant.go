package models

import (
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	DisplayName *string
	StoreName   *string `gorm:"unique"`
	AuthID      string  `gorm:"unique"`
}

func (t *Tenant) HasStore() bool {
	if t.StoreName == nil {
		return false
	}
	return true
}
