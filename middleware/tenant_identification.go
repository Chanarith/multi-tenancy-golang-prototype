package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"core/database"
	"core/models"
)

type TenantIdentificationByPath struct {
}

func NewTenantIdentificationByPath() *TenantIdentificationByPath {
	return &TenantIdentificationByPath{}
}

func (t *TenantIdentificationByPath) Handle(c *gin.Context) {
	tenant := c.Param("tenant")
	var foundTenant models.Tenant
	result := database.GetCentralConnection().Where("store_name = ?", tenant).First(&foundTenant)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		c.Abort()
		return
	}
	Db := database.GetTenantConnection(*foundTenant.StoreName)
	if Db == nil {
		c.Abort()
		return
	}
	c.Set("DB", Db)
	c.Next()
	sqlx, err := Db.DB()
	if err != nil {
		panic(err)
	}
	sqlx.Close()
	centralDB, _ := result.DB()
	centralDB.Close()
}
