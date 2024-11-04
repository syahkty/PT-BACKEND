package delivery

import (
	// "library_api/internal/entity"
	"fmt"
	"library_api/internal/entity"
	"library_api/internal/middleware"
	"library_api/internal/usecase"
	"library_api/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BorrowingHandler struct {
    borrowingUsecase *usecase.BorrowingUsecase
}

func NewBorrowingHandler(r *gin.Engine, uc *usecase.BorrowingUsecase,jwtAuth *auth.JWTAuth) {
    handler := &BorrowingHandler{borrowingUsecase: uc}
    r.GET("/borrowings/user", middleware.AuthMiddleware(jwtAuth),handler.GetBorrowedBooksByUser) // Endpoint untuk user
    r.GET("/borrowings/admin",middleware.AuthMiddleware(jwtAuth), middleware.AdminOnlyMiddleware(),handler.GetAllBorrowings)// Endpoint untuk admin?
    r.POST("/borrowings", middleware.AuthMiddleware(jwtAuth), handler.CreateBorrowing) // Endpoint untuk membuat pinjaman
    r.POST("/borrowings/external", middleware.AuthMiddleware(jwtAuth), handler.CreateBorrowingExternal) // Endpoint untuk membuat pinjaman
    r.PUT("/borrowings/return/:book_id", middleware.AuthMiddleware(jwtAuth), handler.ReturnBook)
}

func (h *BorrowingHandler) GetBorrowedBooksByUser(c *gin.Context) {
	fmt.Println("Handler GetBorrowedBooksByUser called")
    userID, exists := c.Get("id")
	fmt.Println(userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

    borrowings, err := h.borrowingUsecase.GetBorrowedBooksByUserID(c.Request.Context(), userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, borrowings)
}

func (h *BorrowingHandler) GetAllBorrowings(c *gin.Context) {
	fmt.Println("Handler GetAllBorrowings called")
    borrowings, err := h.borrowingUsecase.GetAllBorrowings(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, borrowings)
}

func (h *BorrowingHandler) CreateBorrowing(c *gin.Context) {
	fmt.Println("Handler CreateBorrowing called")
    
	// Ambil User ID dari token (dari middleware)
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Ambil Book ID dari body request
	var requestBody struct {
		BookID uint `json:"book_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Buat record Borrowing
	borrowing := &entity.Borrowing{
		BookID:     requestBody.BookID,
		BookTitle: "",
		UserID:     userID.(uint),
		BorrowedAt: time.Now(),
		Status:     "active",
	}

	// Simpan pinjaman menggunakan usecase
	if err := h.borrowingUsecase.CreateBorrowing(c.Request.Context(), borrowing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Borrowing created successfully"})
}

func (h *BorrowingHandler) CreateBorrowingExternal(c *gin.Context) {
    userID, exists := c.Get("id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    var requestBody struct {
        ExternalBookID uint `json:"external_book_id"`
    }
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Buat record Borrowing untuk buku eksternal
    borrowing := &entity.Borrowing{
        BookID:     requestBody.ExternalBookID, // Simpan ID buku eksternal
		BookTitle :     "",
        UserID:     userID.(uint),
        BorrowedAt: time.Now(),
        Status:     "active",
    }

    if err := h.borrowingUsecase.CreateBorrowingExternal(c.Request.Context(), borrowing); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "External borrowing created successfully"})
}


func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
    bookIDStr := c.Param("book_id")
    bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    // Mengambil userID dari konteks
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
    
    err = h.borrowingUsecase.ReturnBook(c.Request.Context(), uint(bookID),userID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
