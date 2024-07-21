package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	models_tenants "core/models/tenants"
	"core/repositories"
	"core/utils"
)

type ProductController struct{}

func NewProductController() *ProductController {
	return &ProductController{}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsAvailable int     `json:"quantity" binding:"required"`
}

type UpdateProductRequest struct {
	ProductNmae string  `json:"name"`
	Price       float64 `json:"price"`
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	repo := repositories.NewProductRepository(ctx)
	products, err := repo.All(ctx.GetInt("page"), ctx.GetInt("size"))
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Products": products,
	})
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	productRepo := repositories.NewProductRepository(ctx)
	product, err := productRepo.FindByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Product": product,
	})
}

func (c *ProductController) Create(ctx *gin.Context) {
	var request CreateProductRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productRepo := repositories.NewProductRepository(ctx)
	newProduct := models_tenants.Product{
		ProductName: request.Name,
		Price:       request.Price,
		IsAvailable: true,
	}
	createdProduct, err := productRepo.Create(&newProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"product": createdProduct})
}

func (c *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var request UpdateProductRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productRepo := repositories.NewProductRepository(ctx)
	product, err := productRepo.FindByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if request.ProductNmae != "" {
		product.ProductName = request.ProductNmae
	}
	if request.Price != 0 {
		product.Price = request.Price
	}

	_, err = productRepo.Update(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	productRepo := repositories.NewProductRepository(ctx)
	err := productRepo.DeleteByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
