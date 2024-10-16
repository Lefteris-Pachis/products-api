package routes

import (
	"github.com/gin-gonic/gin"
	"products-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductById)
	r.POST("/products", controllers.CreateProduct)
}
