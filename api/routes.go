package api

import (
	"demerzel-badges/internal/handlers"
	"demerzel-badges/internal/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/api/badges/health", handlers.HealthHandler)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api/badges")
	apiRoutes.POST("/badges", handlers.CreateBadgeHandler)
	apiRoutes.GET("/user/badges", middleware.CanViewBadge(), handlers.GetBadgesForUserHandler)
	apiRoutes.POST("/user/badges", middleware.CanAssignBadge(), handlers.AssignBadgeHandler)
	apiRoutes.GET("/user/badges/skill/:skillId", middleware.CanViewBadge(), handlers.GetUserBadgeHandler)
	apiRoutes.GET("/badges/:badge_id", middleware.CanViewBadge(), handlers.GetUserBadgeByIDHandler)

	return r
}
