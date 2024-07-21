package models_tenants

import (
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Products          []Product `gorm:"many2many:invoice_products;"`
	CustomerFirstName string
	CustomerLastName  string
	CustomerEmail     string
	TotalPrice        float64
	TotalCost         float64
}
