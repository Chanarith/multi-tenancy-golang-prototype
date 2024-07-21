package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"core/repositories"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (c *DashboardController) GetSummary(ctx *gin.Context) {
	tenantDB, _ := ctx.Get("DB")
	db := tenantDB.(*gorm.DB)
	dashboardRepo := repositories.DashboardRepository{}
	revenue, totalCost, err := dashboardRepo.GetDashboardSummary(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Revenue": revenue,
			"Expense": totalCost,
			"Profit":  revenue - totalCost,
		})
	}
}
