package api

import (
	"demerzel-badges/internal/handlers"
	"demerzel-badges/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthHandler)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api")
	
	apiRoutes.Use(middleware.UserAuthMiddleware())
	
	apiRoutes.POST("/badges", handlers.CreateBadgeHandler)
	apiRoutes.GET("/badges/:badge_id", handlers.GetUserBadgeByIDHandler)
	apiRoutes.POST("/user/badges", handlers.AssignBadgeHandler)
	apiRoutes.GET("/user/badges/:userId/skill/:skillId", handlers.GetUserBadgeHandler)

	return r
}
