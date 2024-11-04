package delivery

import (
	"library_api/internal/entity"
	"library_api/internal/usecase"
	"library_api/pkg/auth"

	"net/http"
	"strconv"

	"library_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
    bookUsecase *usecase.BookUsecase
}

func NewBookHandler(r *gin.Engine, uc *usecase.BookUsecase, jwtAuth *auth.JWTAuth) {
    handler := &BookHandler{bookUsecase: uc}
    r.POST("/books", middleware.AuthMiddleware(jwtAuth), middleware.AdminOnlyMiddleware(), handler.CreateBook) // Admin
    r.GET("/books", middleware.AuthMiddleware(jwtAuth),handler.GetBooks)                                     // Semua
    r.GET("/books/:id", middleware.AuthMiddleware(jwtAuth),handler.GetBookByID)                             // Semua
    r.GET("/books/external", middleware.AuthMiddleware(jwtAuth), handler.GetBooksFromExternalAPI)
    r.PUT("/books/:id", middleware.AuthMiddleware(jwtAuth), middleware.AdminOnlyMiddleware(), handler.UpdateBook) // Admin
    r.DELETE("/books/:id", middleware.AuthMiddleware(jwtAuth), middleware.AdminOnlyMiddleware(), handler.DeleteBook) // Admin

}
func (h *BookHandler) CreateBook(c *gin.Context) {
    var book entity.Book
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.bookUsecase.CreateBook(c.Request.Context(), &book); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	// Parsing query parameter untuk pagination
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	title := c.Query("title")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	// Panggil usecase untuk mendapatkan buku
	books, total, err := h.bookUsecase.GetBooks(limit, page, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan hasil dengan pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":       books,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Hitung total halaman
	})
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32) 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }
    book, err := h.bookUsecase.GetBookByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }
    c.JSON(http.StatusOK, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }
    var book entity.Book
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    book.ID = uint(id) // Set ID sebagai uint setelah konversi
    if _, err := h.bookUsecase.GetBookByID(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }
    if err := h.bookUsecase.UpdateBook(c.Request.Context(), &book); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32) // Mengonversi string ke uint64
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }
    if err := h.bookUsecase.DeleteBook(c.Request.Context(), uint(id)); err != nil { // Konversi dari uint64 ke uint
        c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}

func (h *BookHandler) GetBooksFromExternalAPI(c *gin.Context) {
    title := c.Query("title")

    // Panggil fungsi pada usecase
    result, err := h.bookUsecase.GetBooksFromExternalAPI(c.Request.Context(), title)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}