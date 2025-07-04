package routes

import (
	"golang-starter-kit/controllers" // Import package controllers untuk mengakses fungsi-fungsi controller
	"golang-starter-kit/middleware"  // Import package middleware untuk mengakses middleware JWT
	"github.com/gin-gonic/gin" 	     // Import framework Gin untuk routing dan handling HTTP requests
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// Public routes
		// Auth
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/logout", middleware.JWTAuth(), controllers.Logout)

		// Secret
		secret := api.Group("/secret")
		{
			// Black List
			secret.POST("/get-black-list", controllers.GetBlacklistTokens)
			secret.POST("/clear-black-list", controllers.ClearBlacklistTokens)
		}

		// User
		user := api.Group("/user")
		{
			user.GET("/", middleware.JWTAuth(), controllers.GetUsers)
			user.GET("/:id", middleware.JWTAuth(), controllers.GetUserByID)
			user.POST("/", middleware.JWTAuth(), controllers.CreateUser)
			user.PUT("/:id", middleware.JWTAuth(), controllers.UpdateUser)
			user.DELETE("/:id", middleware.JWTAuth(), controllers.DeleteUser)
		}

		// Role
		role := api.Group("/role")
		{
			role.GET("/", middleware.JWTAuth(), controllers.GetRoles)
			role.GET("/:id", middleware.JWTAuth(), controllers.GetRoleByID)
			role.POST("/", middleware.JWTAuth(), controllers.CreateRole)
			role.PUT("/:id", middleware.JWTAuth(), controllers.UpdateRole)
			role.DELETE("/:id", middleware.JWTAuth(), controllers.DeleteRole)
		}
	}

	return r
}
