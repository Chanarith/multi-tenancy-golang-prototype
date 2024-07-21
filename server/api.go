package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"

	"core/controllers"
	"core/middleware"
	"core/repositories"
)

func Run() {
	server := gin.Default()
	// start supertokens
	server.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(os.Getenv("SUPERTOKENS_ALLOWED_ORIGINS"), ","),
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: append([]string{"content-type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))
	server.Use(func(ctx *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	})
	// end supertokens
	api := server.Group("/api/v1")
	{
		me := api.Group("/me")
		{
			me.GET("/",
				middleware.NewForceCentralConnection().Handle,
				middleware.NewVerifySession(nil).Handle,
				func(ctx *gin.Context) {
					sessionContainer := session.GetSessionFromRequestContext(ctx.Request.Context())
					userID := sessionContainer.GetUserID()
					userInfo, errGetUserByIdST := emailpassword.GetUserByID(userID)
					if errGetUserByIdST != nil {
						return
					}
					repo := repositories.NewTenantRepository(ctx)
					tenantProfile, errGetUserById := repo.GetByID(userID)
					if errGetUserById != nil {
						return
					}
					ctx.JSON(http.StatusOK, gin.H{
						"Email":       userInfo.Email,
						"StoreName":   tenantProfile.StoreName,
						"DisplayName": tenantProfile.DisplayName,
					})
				})
		}
		landlordApi := api.Group("/tenants", middleware.NewForceCentralConnection().Handle)
		{
			landlordApi.GET("/", middleware.DefaultPaginator().Handle, controllers.NewTenantController().GetAll)
			landlordApi.GET("/:id", controllers.NewTenantController().GetByID)
			landlordApi.POST("/store", middleware.NewVerifySession(nil).Handle, controllers.NewTenantController().CreateStore)
		}
		tenantApi := api.Group("/:tenant", middleware.NewTenantIdentificationByPath().Handle)
		{
			summaryApi := tenantApi.Group("/")
			{
				summaryApi.GET("/", controllers.NewDashboardController().GetSummary)
			}
			productApi := tenantApi.Group("/products")
			{
				productApi.POST("/", middleware.NewVerifySession(nil).Handle, controllers.NewProductController().Create)
				productApi.GET("/", middleware.DefaultPaginator().Handle, controllers.NewProductController().GetAll)
				productApi.GET("/:id", controllers.NewProductController().GetByID)
				productApi.PATCH("/:id", middleware.NewVerifySession(nil).Handle, controllers.NewProductController().Update)
				productApi.DELETE("/:id", middleware.NewVerifySession(nil).Handle, controllers.NewProductController().Delete)
			}
			invoiceApi := tenantApi.Group("/invoices")
			{
				invoiceApi.POST("/", controllers.NewInvoiceController().Create)
				invoiceApi.GET("/", middleware.DefaultPaginator().Handle, controllers.NewInvoiceController().GetAll)
				invoiceApi.GET("/:id", controllers.NewInvoiceController().GetByID)
				invoiceApi.PATCH("/:id", controllers.NewInvoiceController().Update)
				invoiceApi.DELETE("/:id", controllers.NewInvoiceController().Delete)
			}
		}
	}
	server.Run()
}
