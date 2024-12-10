package routes

import (
	"log"
	"net/http"
	"project/helper"
	"project/infra"
	"sync"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoutes(ctx infra.ServiceContext) *http.Server {
	r := gin.Default()

	r.Use(ctx.Middleware.Logger())
	r.POST("/login", ctx.Ctl.AuthHandler.Login)
	r.POST("/register", ctx.Ctl.UserHandler.Registration)
	r.GET("/users", ctx.Ctl.UserHandler.All)
	r.POST("/password-reset", ctx.Ctl.PasswordResetHandler.Create)

	category := r.Group("/category")
	{
		category.GET("/", ctx.Ctl.Category.ShowAllCategory)
		category.POST("/", ctx.Ctl.Category.CreateCategory)
		category.DELETE("/:id", ctx.Middleware.OnlyAdmin(), ctx.Ctl.Category.DeleteCategory)
		category.GET("/:id", ctx.Ctl.Category.GetCategoryByID)
		category.PUT("/:id", ctx.Ctl.Category.UpdateCategory)
	}

	banner := r.Group("/banner")
	{
		banner.GET("/", ctx.Ctl.Banner.GetAll)
		banner.POST("/", ctx.Ctl.Banner.Create)
		banner.GET("/:id", ctx.Ctl.Banner.GetById)
		banner.PUT("/:id", ctx.Ctl.Banner.Edit)
		banner.DELETE("/:id", ctx.Middleware.OnlyAdmin(), ctx.Ctl.Banner.Delete)
	}

	products := r.Group("/products")
	{
		products.GET("/", ctx.Ctl.Product.ShowAllProduct)
		products.POST("/", ctx.Ctl.Product.CreateProduct)
		products.GET("/:id", ctx.Ctl.Product.GetProductByID)
		products.DELETE("/:id", ctx.Middleware.OnlyAdmin(), ctx.Ctl.Product.DeleteProduct)
		products.PUT("/:id", ctx.Ctl.Product.UpdateProduct)
	}

	order := r.Group("/orders")
	{
		order.GET("/", ctx.Ctl.OrderHandler.All)
		order.GET("/:id", ctx.Ctl.OrderHandler.Get)
		order.PUT("/:id", ctx.Ctl.OrderHandler.Update)
	}

	dashboard := r.Group("dashboard")
	{
		dashboard.GET("/earning", ctx.Ctl.Dashboard.GetEarningDashboard)
		dashboard.GET("/summary", ctx.Ctl.Dashboard.GetSummary)
		dashboard.GET("/bestSeller", ctx.Ctl.Dashboard.GetBestSeller)
		dashboard.GET("/revenue", ctx.Ctl.Dashboard.GetMonthlyRevenue)
	}

	stock := r.Group("/stock")
	{
		stock.GET("/:productVariantId", ctx.Ctl.Stock.GetDetails)
		stock.PUT("/:productVariantId", ctx.Ctl.Stock.Edit)
		stock.DELETE("/:id", ctx.Middleware.OnlyAdmin(), ctx.Ctl.Stock.Delete)
	}

	promotion := r.Group("/promotion")
	{
		promotion.GET("/", ctx.Ctl.Promotion.GetAll)
		promotion.GET("/:id", ctx.Ctl.Promotion.GetById)
		promotion.POST("/", ctx.Ctl.Promotion.Create)
		promotion.DELETE("/:id", ctx.Middleware.OnlyAdmin(), ctx.Ctl.Promotion.Delete)

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/cdn-upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["images[]"]

		var wg sync.WaitGroup
		responses, _ := helper.Upload(&wg, files)
		log.Println(responses)
	})

	return &http.Server{
		Addr:    ctx.Cfg.ServerPort,
		Handler: r.Handler(),
	}
}
