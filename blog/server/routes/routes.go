package routes

import (
	"blog-server/controllers"
	"blog-server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		api.POST("/login", controllers.Login)

		api.GET("/articles", controllers.GetPublishedArticles)
		api.GET("/articles/:id", controllers.GetArticle)
		api.GET("/articles/:id/comments", controllers.GetArticleComments)
		api.GET("/captcha", controllers.GenerateCaptcha)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/tags", controllers.GetTags)

		api.POST("/comments", controllers.CreateComment)
		api.POST("/comments/:id/like", controllers.LikeComment)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			admin.GET("/articles", controllers.GetAdminArticles)
			admin.POST("/articles", controllers.CreateArticle)
			admin.PUT("/articles/:id", controllers.UpdateArticle)
			admin.DELETE("/articles/:id", controllers.DeleteArticle)
			admin.PUT("/password", controllers.ChangePassword)

			admin.POST("/upload", controllers.UploadFile)

			admin.GET("/categories", controllers.GetCategories)
			admin.POST("/categories", controllers.CreateCategory)
			admin.DELETE("/categories/:id", controllers.DeleteCategory)

			admin.GET("/tags", controllers.GetTags)
			admin.POST("/tags", controllers.CreateTag)
			admin.DELETE("/tags/:id", controllers.DeleteTag)

			admin.GET("/comments", controllers.GetAdminComments)
			admin.PUT("/comments/:id/status", controllers.UpdateCommentStatus)
			admin.DELETE("/comments/:id", controllers.DeleteComment)
			admin.POST("/comments/batch-delete", controllers.DeleteCommentsBatch)

			admin.GET("/settings", controllers.GetSystemSettings)
			admin.PUT("/settings", controllers.UpdateSystemSettings)
		}
	}
}
