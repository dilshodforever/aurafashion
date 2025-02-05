package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/cors" // Import the CORS package
	"github.com/gin-gonic/gin"
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
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis)

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


	user := v1.Group("/user")
	{
		user.POST("/user/", handlerV1.CreateUser)
		//user.GET("/user/list", handlerV1.GetUsers)
		user.GET("/user/:id", handlerV1.GetUser)
		user.PUT("/user/", handlerV1.UpdateUser)
		user.DELETE("/user/:id", handlerV1.DeleteUser)
	}


	
	{
		v1.POST("/auth/logout", handlerV1.Logout)
		v1.POST("/auth/register", handlerV1.Register)
		v1.POST("/auth/verify-email", handlerV1.VerifyEmail)
		v1.POST("/auth/login", handlerV1.Login)
	}

	product := v1.Group("/product")
	{
		product.POST("/product/", handlerV1.CreateProduct)
		product.PUT("/product", handlerV1.UpdateProduct)
		product.DELETE("/product/:id", handlerV1.DeleteProduct)
		product.GET("/product/list", handlerV1.ListProducts)
		//product.GET("/product/:id", handlerV1.GetProduct)
		product.POST("/product/picture", handlerV1.AddPicture)
		product.DELETE("/product/picture", handlerV1.DeletePicture)
	}

	basket := v1.Group("/basket")
	{
		basket.POST("/basket/item", handlerV1.AddBasketItem)
		basket.DELETE("/basket", handlerV1.DeleteBasket)
		basket.DELETE("/basket/item", handlerV1.DeleteBasketItem)
		basket.GET("/basket", handlerV1.GetBasket)
	}

	order := v1.Group("/order")
	{
		order.POST("/order", handlerV1.CreateOrder)
		order.PUT("/order", handlerV1.UpdateOrder)
		order.DELETE("/order/:id", handlerV1.DeleteOrder)
		order.GET("/order/list", handlerV1.ListOrders)
		//order.GET("/order/:id", handlerV1.GetOrder)
		order.GET("/order/products", handlerV1.SeeOrderProducts)
	}
	post := v1.Group("/post")
	{
		post.POST("/post", handlerV1.CreatePost)
		post.PUT("/post", handlerV1.UpdatePost)
		post.DELETE("/post/:id", handlerV1.DeletePost)
		post.GET("/post/list", handlerV1.ListPosts)
		//post.GET("/post/:id", handlerV1.GetPost)
		post.POST("/post/picture", handlerV1.AddPostPicture)
		post.DELETE("/post/picture", handlerV1.DeletePostPicture)
	}
	category := v1.Group("/category")
	{
		category.POST("/category", handlerV1.CreateCategory)
		category.GET("/category/:id", handlerV1.GetCategory)
		category.GET("/category/list", handlerV1.GetCategories)
		category.PUT("/category/:id", handlerV1.UpdateCategory)
		category.DELETE("/category/:id", handlerV1.DeleteCategory)
	}


}