package routes

import (
	"github.com/gin-gonic/gin"
	"clientgo/controllers"
)

func RegisterRouter(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/generate", controllers.GenerateRecordsHandler)
	router.GET("/count", controllers.CountRecordsHandler)

}
