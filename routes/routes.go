package routes

import (
	"github.com/gin-gonic/gin"
	"products-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/products", controllers.CreateProduct)
}
