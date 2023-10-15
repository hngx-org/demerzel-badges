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

	r.GET("/health", handlers.HealthHandler)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api")

	apiRoutes.POST("/badges", handlers.CreateBadgeHandler)

	apiRoutes.POST("/user/badges", func(c *gin.Context) {
		userID, _ := c.Get(middleware.UserIDKey)

		targetUserID := c.PostForm("user_id")

		if userID != targetUserID {
			response.Error(c, http.StatusForbidden, "Forbidden", map[string]interface{}{
				"error": "You can only assign badges to yourself",
			})
			c.Abort()
			return
		}
		handlers.AssignBadgeHandler(c)
	})  	
	
	apiRoutes.GET("/user/badges/:userId/skill/:skillId", handlers.GetUserBadgeHandler)
	apiRoutes.GET("/badges/:badge_id", handlers.GetUserBadgeByIDHandler)

	return r
}
