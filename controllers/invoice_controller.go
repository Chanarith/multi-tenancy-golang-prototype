package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	models_tenants "core/models/tenants"
	"core/repositories"
	"core/utils"
)

type InvoiceController struct{}

func NewInvoiceController() *InvoiceController {
	return &InvoiceController{}
}

type productRequest struct {
	ID uint `json:"id" binding:"required"`
}

type CreateInvoiceRequest struct {
	Products          []*productRequest `json:"products" binding:"required,min=1,dive,product_exist"`
	CustomerFirstName string            `json:"customer_first_name" binding:"required"`
	CustomerLastName  string            `json:"customer_last_name" binding:"required"`
	CustomerEmail     string            `json:"customer_email" binding:"required,email"`
	TotalDiscount     float64           `json:"total_discount" binding:"required"`
	VAT               float64           `json:"vat" binding:"required"`
}

func productExist(ctx *gin.Context) validator.Func {
	return func(fl validator.FieldLevel) bool {
		intf := fl.Field().Interface().(productRequest)
		repo := repositories.NewProductRepository(ctx)
		if _, err := repo.FindByID(intf.ID); err != nil {
			return false
		}
		return true
	}
}

func (c *InvoiceController) GetAll(ctx *gin.Context) {
	repo := repositories.NewInvoiceRepository(ctx)
	invoices, err := repo.All()
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Invoices": invoices,
	})
}

func (c *InvoiceController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	invoiceRepo := repositories.NewInvoiceRepository(ctx)
	invoice, err := invoiceRepo.FindByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invoice not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Invoice": invoice,
	})
}

func (c *InvoiceController) Create(ctx *gin.Context) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var f validator.Func = productExist(ctx)
		v.RegisterValidation("product_exist", f)
	}
	var request CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceRepo := repositories.NewInvoiceRepository(ctx)
	var productIds []uint
	for _, product := range request.Products {
		productIds = append(productIds, product.ID)
	}
	productRepo := repositories.NewProductRepository(ctx)
	products, err := productRepo.GetByIDs(productIds)
	var convertedProducts []models_tenants.Product
	for _, product := range products {
		convertedProducts = append(convertedProducts, *product)
	}
	var totalPrice float64
	for _, product := range convertedProducts {
		totalPrice += product.Price
	}
	var totalCost float64
	for _, product := range convertedProducts {
		totalPrice += product.CostOfGoodSold
	}
	newInvoice := models_tenants.Invoice{
		Products:          convertedProducts,
		CustomerFirstName: request.CustomerFirstName,
		CustomerLastName:  request.CustomerLastName,
		CustomerEmail:     request.CustomerEmail,
		TotalPrice:        totalPrice,
		TotalCost:         totalCost,
	}
	createdInvoice, _ := invoiceRepo.Create(&newInvoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"invoice": createdInvoice})
}

func (c *InvoiceController) Update(ctx *gin.Context) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var f validator.Func = productExist(ctx)
		v.RegisterValidation("product_exist", f)
	}
	id := ctx.Param("id")
	var request CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceRepo := repositories.NewInvoiceRepository(ctx)
	invoice, err := invoiceRepo.FindByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	var productIds []uint
	for _, product := range request.Products {
		productIds = append(productIds, product.ID)
	}
	productRepo := repositories.NewProductRepository(ctx)
	products, err := productRepo.GetByIDs(productIds)
	var convertedProducts []models_tenants.Product
	for _, product := range products {
		convertedProducts = append(convertedProducts, *product)
	}
	var totalPrice float64
	for _, product := range convertedProducts {
		totalPrice += product.Price
	}
	var totalCost float64
	for _, product := range convertedProducts {
		totalPrice += product.CostOfGoodSold
	}
	invoice.Products = convertedProducts
	invoice.CustomerFirstName = request.CustomerFirstName
	invoice.CustomerLastName = request.CustomerLastName
	invoice.CustomerEmail = request.CustomerEmail
	invoice.TotalPrice = totalPrice
	invoice.TotalCost = totalCost
	updatedInvoice, _ := invoiceRepo.Update(invoice)
	invoiceRepo.AttachWithReplaceProducts(updatedInvoice, products)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully"})
}

func (c *InvoiceController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	invoiceRepo := repositories.InvoiceRepository{}
	err := invoiceRepo.DeleteByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
