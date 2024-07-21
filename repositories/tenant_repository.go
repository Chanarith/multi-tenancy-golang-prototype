package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"core/database"
	"core/models"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(ctx *gin.Context) *TenantRepository {
	tenantDB, _ := ctx.Get("DB")
	db := tenantDB.(*gorm.DB)
	return &TenantRepository{db: db}
}

func NewTenantRepositoryFromDB(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) All(page int, pageSize int) ([]models.Tenant, error) {
	var tenants []models.Tenant
	results := r.db.Scopes(database.ScopePaginate(page, pageSize)).Find(&tenants)
	if results.Error != nil {
		println(results.Error)
		return nil, results.Error
	}
	return tenants, nil
}

func (r *TenantRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.Tenant{}).Count(&count)
	if result.Error != nil {
		println(result.Error)
		return 0, result.Error
	}
	return count, nil
}

func (r *TenantRepository) First() (models.Tenant, error) {
	var tenant models.Tenant
	result := r.db.First(&tenant)
	if result.Error != nil {
		println(result.Error)
		return tenant, result.Error
	}
	return tenant, nil
}

func (r *TenantRepository) GetByID(id string) (models.Tenant, error) {
	var tenant models.Tenant
	result := r.db.Where("auth_id = ?", id).First(&tenant)
	if result.Error != nil {
		println(result.Error)
		return tenant, result.Error
	}
	return tenant, nil
}

func (r *TenantRepository) Create(newTenant models.Tenant) (models.Tenant, error) {
	result := r.db.Create(&newTenant)
	if result.Error != nil {
		println(result.Error)
		return models.Tenant{}, result.Error
	}
	return newTenant, nil
}

func (r *TenantRepository) Update(tenant models.Tenant) (models.Tenant, error) {
	result := r.db.Omit("AuthID").Where("auth_id", tenant.AuthID).Updates(tenant)
	if result.Error != nil {
		println(result.Error)
		return models.Tenant{}, result.Error
	}
	return tenant, nil
}

func (r *TenantRepository) DeleteByID(id uint) error {
	var tenant models.Tenant
	result := r.db.Delete(&tenant, id)
	if result.Error != nil {
		println(result.Error)
		return result.Error
	}
	return nil
}
