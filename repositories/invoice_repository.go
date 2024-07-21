package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	models_tenants "core/models/tenants"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(ctx *gin.Context) *InvoiceRepository {
	tenantDB, _ := ctx.Get("DB")
	db := tenantDB.(*gorm.DB)
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) All() ([]*models_tenants.Invoice, error) {
	var invoices []*models_tenants.Invoice
	results := r.db.Preload("Products").Find(&invoices)
	if results.Error != nil {
		println(results.Error)
		return nil, results.Error
	}
	return invoices, nil
}

func (r *InvoiceRepository) FindByID(id uint) (*models_tenants.Invoice, error) {
	var invoice models_tenants.Invoice
	result := r.db.Preload("Products").First(&invoice, id)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return &invoice, nil
}

func (r *InvoiceRepository) Create(newInvoice *models_tenants.Invoice) (*models_tenants.Invoice, error) {
	result := r.db.Create(newInvoice)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return newInvoice, nil
}

func (r *InvoiceRepository) AttachProduct(invoice *models_tenants.Invoice, product *models_tenants.Product) (*models_tenants.Product, error) {
	db := r.db
	err := db.Model(&invoice).Association("Products").Append(&product)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return product, nil
}

func (r *InvoiceRepository) AttachProducts(invoice *models_tenants.Invoice, products []models_tenants.Product) ([]models_tenants.Product, error) {
	db := r.db
	// Start a new transaction
	tx := db.Begin()
	for _, product := range products {
		err := tx.Model(&invoice).Association("Products").Append(&product)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
			println(err.Error())
			return nil, err
		}
	}
	// Commit the transaction if all attachments were successful
	tx.Commit()
	return products, nil
}

func (r *InvoiceRepository) AttachWithReplaceProducts(invoice *models_tenants.Invoice, products []*models_tenants.Product) ([]*models_tenants.Product, error) {
	db := r.db
	// Start a new transaction
	tx := db.Begin()
	for _, product := range products {
		err := tx.Model(&invoice).Association("Products").Replace(product)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
			println(err.Error())
			return nil, err
		}
	}
	// Commit the transaction if all attachments were successful
	tx.Commit()
	return products, nil
}

func (r *InvoiceRepository) Update(invoice *models_tenants.Invoice) (*models_tenants.Invoice, error) {
	result := r.db.Save(invoice)
	if result.Error != nil {
		println(result.Error)
		return nil, result.Error
	}
	return invoice, nil
}

func (r *InvoiceRepository) DeleteByID(id uint) error {
	var invoice models_tenants.Invoice
	result := r.db.Delete(&invoice, id)
	if result.Error != nil {
		println(result.Error)
		return result.Error
	}
	return nil
}
