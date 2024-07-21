package repositories

import (
	"fmt"

	"gorm.io/gorm"

	models_tenants "core/models/tenants"
)

type DashboardRepository struct{}

func (r *DashboardRepository) GetDashboardSummary(db *gorm.DB) (float64, float64, error) {
	var revenue float64
	if err := db.Model(&models_tenants.Invoice{}).Select("COALESCE(SUM(total_price), 0) as revenue").Scan(&revenue).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to calculate revenue: %v", err)
	}
	var totalCost float64
	if err := db.Model(&models_tenants.Invoice{}).Select("COALESCE(SUM(total_cost), 0) as total_cost").Scan(&totalCost).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to calculate total cost: %v", err)
	}
	return revenue, totalCost, nil
}
