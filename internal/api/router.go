package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the API router
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		// Scheduler operations
		scheduler := api.Group("/scheduler")
		{
			scheduler.POST("/start", handler.StartScheduler)
			scheduler.POST("/stop", handler.StopScheduler)
			scheduler.GET("/status", handler.GetSchedulerStatus)
		}

		// Message operations
		messages := api.Group("/messages")
		{
			messages.GET("", handler.GetSentMessages)
		}
	}

	return r
}
