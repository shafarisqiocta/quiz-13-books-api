package main

import (
	"log"
	"os"
	"quiz-13-books-api/config"
	"quiz-13-books-api/handlers"
	"quiz-13-books-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config.ConnectDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server running",
		})
	})

	// group protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.BasicAuthMiddleware())

	// ===== CATEGORY =====
	authorized.GET("/api/categories", handlers.GetCategories)
	authorized.POST("/api/categories", handlers.CreateCategory)
	authorized.GET("/api/categories/:id", handlers.GetCategoryByID)
	authorized.DELETE("/api/categories/:id", handlers.DeleteCategory)
	authorized.GET("/api/categories/:id/books", handlers.GetBooksByCategory)
	authorized.PUT("/api/categories/:id", handlers.UpdateCategory)

	// ===== BOOK =====
	authorized.GET("/api/books", handlers.GetBooks)
	authorized.POST("/api/books", handlers.CreateBook)
	authorized.GET("/api/books/:id", handlers.GetBookByID)
	authorized.DELETE("/api/books/:id", handlers.DeleteBook)
	authorized.PUT("/api/books/:id", handlers.UpdateBook)

	// ✅ UBAH BAGIAN INI - Port dinamis untuk Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default untuk local
	}
	r.Run(":" + port)
}
