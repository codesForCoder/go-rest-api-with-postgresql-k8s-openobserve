package routes

import (
	"go-rest-api/controllers"
	"go-rest-api/middlewares"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//Register Tarace Provider
	r.Use(otelgin.Middleware("go-crud"))
	// Register the IPLogger middleware
	r.Use(middlewares.IPLogger())

	r.GET("/products", controllers.FindProducts)
	r.POST("/products", controllers.CreateProduct)
	r.GET("/products/:id", controllers.FindProduct)
	r.PUT("/products/:id", controllers.UpdateProduct)
	r.DELETE("/products/:id", controllers.DeleteProduct)

	return r
}
