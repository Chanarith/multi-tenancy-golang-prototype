package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"core/database"
	models_tenants "core/models/tenants"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(ctx *gin.Context) *ProductRepository {
	DB, _ := ctx.Get("DB")
	db := DB.(*gorm.DB)
	return &ProductRepository{db}
}

func (r *ProductRepository) All(page int, pageSize int) ([]*models_tenants.Product, error) {
	var products []*models_tenants.Product
	results := r.db.Scopes(database.ScopePaginate(page, pageSize)).Find(&products)
	if results.Error != nil {
		println(results.Error)
		return nil, results.Error
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id uint) (*models_tenants.Product, error) {
	var product models_tenants.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDs(ids []uint) ([]*models_tenants.Product, error) {
	var products []*models_tenants.Product
	result := r.db.Where("id IN ?", ids).Find(&products)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return products, nil
}

func (r *ProductRepository) Create(newProduct *models_tenants.Product) (*models_tenants.Product, error) {
	result := r.db.Create(newProduct)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return newProduct, nil
}

func (r *ProductRepository) Update(product *models_tenants.Product) (*models_tenants.Product, error) {
	result := r.db.Save(product)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return product, nil
}

func (r *ProductRepository) DeleteByID(id uint) error {
	var product models_tenants.Product
	result := r.db.Delete(&product, id)
	if result.Error != nil {
		println(result.Error)
		return result.Error
	}
	return nil
}
