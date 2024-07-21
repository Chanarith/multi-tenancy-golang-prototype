package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/supertokens/supertokens-golang/recipe/session"

	"core/jobs"
	"core/models"
	"core/repositories"
	"core/utils"
)

type TenantController struct {
}

func NewTenantController() *TenantController {
	return &TenantController{}
}

type CreateStoreRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTenantRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
}

func (c *TenantController) GetAll(ctx *gin.Context) {
	repo := repositories.NewTenantRepository(ctx)
	tenants, err := repo.All(ctx.GetInt("page"), ctx.GetInt("size"))
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}
	count, _ := repo.Count()
	ctx.JSON(http.StatusOK, gin.H{
		"Tenants":   tenants,
		"TotalSize": count,
		"Page":      ctx.GetInt("page"),
		"Size":      ctx.GetInt("size"),
	})
}

func (c *TenantController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctx.ShouldBindUri(struct {
		id string `uri:"id" binding:"required"`
	}{id: id}); err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}
	repo := repositories.NewTenantRepository(ctx)
	tenant, err := repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Tenant not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

func (c *TenantController) Create(ctx *gin.Context) {
	var newTenant models.Tenant
	if err := ctx.ShouldBindJSON(&newTenant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := repositories.NewTenantRepository(ctx)
	createdTenant, err := repo.Create(newTenant)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"tenant": createdTenant})
}

func (c *TenantController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var request UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := repositories.NewTenantRepository(ctx)
	tenant, err := repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	newTenant := models.Tenant{DisplayName: &request.DisplayName}
	newTenant.ID = tenant.ID
	_, err = repo.Update(newTenant)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tenant"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
}

func (c *TenantController) CreateStore(ctx *gin.Context) {
	var request CreateStoreRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessionContainer := session.GetSessionFromRequestContext(ctx.Request.Context())
	authID := sessionContainer.GetUserID()
	repo := repositories.NewTenantRepository(ctx)
	if _, err := repo.GetByID(authID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	tenant := models.Tenant{
		AuthID:    authID,
		StoreName: &request.Name,
	}
	if _, err := repo.Update(tenant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := jobs.Perform("onboard_tenant", work.Q{
		"name": &tenant.StoreName,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": true})
}

func (c *TenantController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	repo := repositories.NewTenantRepository(ctx)
	err := repo.DeleteByID(utils.StrToUnint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tenant"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}
