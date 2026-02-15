package handlers

import (
	"database/sql"
	"net/http"
	"quiz-13-books-api/config"
	"quiz-13-books-api/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		var modifiedAt sql.NullTime
		var modifiedBy sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.CreatedBy,
			&modifiedAt,
			&modifiedBy,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// assign manual
		if modifiedAt.Valid {
			category.ModifiedAt = &modifiedAt.Time
		}

		if modifiedBy.Valid {
			category.ModifiedBy = &modifiedBy.String
		}

		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var input struct {
		Name string `json:"Name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `
		INSERT INTO categories (name, created_at, created_by)
		VALUES ($1, NOW(), 'admin')
		RETURNING id
	`
	var id int
	err := config.DB.QueryRow(query, input.Name).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created",
		"id":      id,
	})
}
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT id, name, created_at, created_by, modified_at, modified_by
		FROM categories
		WHERE id = $1
	`

	var category models.Category

	err := config.DB.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.CreatedBy,
		&category.ModifiedAt,
		&category.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	result, err := config.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted",
	})
}
func GetBooksByCategory(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT b.id, b.title, b.description, b.image_url,
		       b.release_year, b.price, b.total_page,
		       b.thickness, b.category_id
		FROM books b
		WHERE b.category_id = $1
	`

	rows, err := config.DB.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []map[string]interface{}

	for rows.Next() {
		var (
			id          int
			title       string
			description string
			imageURL    string
			releaseYear int
			price       int
			totalPage   int
			thickness   string
			categoryID  int
		)

		err := rows.Scan(
			&id,
			&title,
			&description,
			&imageURL,
			&releaseYear,
			&price,
			&totalPage,
			&thickness,
			&categoryID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		book := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"image_url":    imageURL,
			"release_year": releaseYear,
			"price":        price,
			"total_page":   totalPage,
			"thickness":    thickness,
			"category_id":  categoryID,
		}

		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `UPDATE categories SET name=$1 WHERE id=$2`

	result, err := config.DB.Exec(query, input.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated",
	})
}
