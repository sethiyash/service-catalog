package router

import (
	"service-catalog/internal/handlers"
	"service-catalog/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(client *mongo.Client) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", client)
	})

	r.Use(middleware.AuthMiddleware())

	r.GET("/services", handlers.ListServices)
	r.GET("/services/:id", handlers.GetService)
	r.POST("/services", handlers.CreateService)
	r.PUT("/services/:id", handlers.UpdateService)
	r.DELETE("/services/:id", handlers.DeleteService)

	return r
}
