package models_tenants

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName    string
	Price          float64
	CostOfGoodSold float64
	IsAvailable    bool
}
