package delivery

import (
    "net/http"
    "library_api/internal/entity"
    "library_api/internal/usecase"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userUsecase *usecase.UserUsecase
}

func NewUserHandler(r *gin.Engine, uc *usecase.UserUsecase) {
    handler := &UserHandler{userUsecase: uc}
    r.POST("/register", handler.RegisterUser)
    r.POST("/login", handler.LoginUser)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
    var user entity.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userUsecase.RegisterUser(c.Request.Context(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.userUsecase.LoginUser(c.Request.Context(), loginData.Username, loginData.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
