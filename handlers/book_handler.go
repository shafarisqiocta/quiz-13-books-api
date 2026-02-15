package handlers

import (
	"net/http"
	"quiz-13-books-api/config"
	"quiz-13-books-api/models"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url,
		       release_year, price, total_page,
		       thickness, category_id,
		       created_at, created_by
		FROM books
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var book models.Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.ImageURL,
			&book.ReleaseYear,
			&book.Price,
			&book.TotalPage,
			&book.Thickness,
			&book.CategoryID,
			&book.CreatedAt,
			&book.CreatedBy,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}
func CreateBook(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		ReleaseYear int    `json:"release_year"`
		Price       int    `json:"price"`
		TotalPage   int    `json:"total_page"`
		CategoryID  int    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi tahun rilis
	if input.ReleaseYear < 1980 || input.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "release_year must be between 1980 and 2024",
		})
		return
	}

	// validasi thickness
	thickness := "tipis"
	if input.TotalPage > 100 {
		thickness = "tebal"
	}

	query := `
		INSERT INTO books (
			title, description, image_url,
			release_year, price, total_page,
			thickness, category_id,
			created_at, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),'admin')
		RETURNING id
	`

	var id int
	err := config.DB.QueryRow(query,
		input.Title,
		input.Description,
		input.ImageURL,
		input.ReleaseYear,
		input.Price,
		input.TotalPage,
		thickness,
		input.CategoryID,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created",
		"id":      id,
	})
}
func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT id, title, description, image_url,
		       release_year, price, total_page,
		       thickness, category_id,
		       created_at, created_by
		FROM books
		WHERE id = $1
	`

	var book models.Book

	err := config.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.ImageURL,
		&book.ReleaseYear,
		&book.Price,
		&book.TotalPage,
		&book.Thickness,
		&book.CategoryID,
		&book.CreatedAt,
		&book.CreatedBy,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	result, err := config.DB.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted",
	})
}
func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var input models.Book

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// logic thickness
	thickness := "tipis"
	if input.TotalPage > 100 {
		thickness = "tebal"
	}

	query := `
		UPDATE books
		SET title=$1,
		    description=$2,
		    image_url=$3,
		    release_year=$4,
		    price=$5,
		    total_page=$6,
		    thickness=$7,
		    category_id=$8
		WHERE id=$9
	`

	result, err := config.DB.Exec(query,
		input.Title,
		input.Description,
		input.ImageURL,
		input.ReleaseYear,
		input.Price,
		input.TotalPage,
		thickness,
		input.CategoryID,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated",
	})
}
