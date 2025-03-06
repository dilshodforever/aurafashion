package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/cors" // Import the CORS package
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"aura-fashion/config"
	_ "aura-fashion/docs"
	"aura-fashion/internal/controller/http/v1/handler"
	"aura-fashion/internal/usecase"
	"aura-fashion/pkg/logger"

	rediscache "github.com/golanguzb70/redis-cache"
)

// NewRouter -.
// Swagger spec:
// @title       The Muallimah API
// @description This is a sample server The Muallimah server.
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache, MinIO *minio.Client) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis,MinIO)

	// Initialize Casbin enforcer

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend domenini yozish
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e))

	// Swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1")

	{
		//v1.GET("/user/list", handlerV1.GetUsers)
		v1.GET("/user/:id", handlerV1.GetUser)
		v1.PUT("/user/", handlerV1.UpdateUser)
		v1.DELETE("/user/:id", handlerV1.DeleteUser)
	}

	{
		v1.POST("/auth/register", handlerV1.Register)
		v1.POST("/auth/verify-email", handlerV1.VerifyEmail)
		v1.POST("/auth/login", handlerV1.Login)
	}

	{
		v1.POST("/product/", handlerV1.CreateProduct)
		v1.PUT("/product/:id", handlerV1.UpdateProduct)
		v1.DELETE("/product/:id", handlerV1.DeleteProduct)
		v1.GET("/product/list", handlerV1.ListProducts)
		//product.GET("/product/:id", handlerV1.GetProduct)
		v1.POST("/product/picture", handlerV1.AddPicture)
		v1.DELETE("/product/picture", handlerV1.DeletePicture)
	}

	{
		v1.POST("/basket/item", handlerV1.AddBasketItem)
		v1.DELETE("/basket/", handlerV1.DeleteBasket)
		v1.DELETE("/basket/item", handlerV1.DeleteBasketItem)
		v1.GET("/basket/get", handlerV1.GetBasket)
	}

	{
		v1.POST("/order", handlerV1.CreateOrder)
		v1.PUT("/order", handlerV1.UpdateOrder)
		v1.DELETE("/order/:id", handlerV1.DeleteOrder)
		v1.GET("/order/list", handlerV1.ListOrders)
		//order.GET("/order/:id", handlerV1.GetOrder)
		v1.GET("/order/products", handlerV1.SeeOrderProducts)
	}

	{
		v1.POST("/category", handlerV1.CreateCategory)
		v1.GET("/category/:id", handlerV1.GetCategory)
		v1.GET("/category/list", handlerV1.GetCategories)
		v1.PUT("/category/:id", handlerV1.UpdateCategory)
		v1.DELETE("/category/:id", handlerV1.DeleteCategory)
	}
	v1.POST("/minio/media", handlerV1.Media)
}
