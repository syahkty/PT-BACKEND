package main

import (
    "log"
    "library_api/configs"
    "library_api/internal/delivery"
    "library_api/internal/repository"
    "library_api/internal/usecase"
    "library_api/pkg/auth"
    "github.com/gin-gonic/gin"
)

func main() {
    // Connect to the database
    db, err := configs.ConnectDatabase()
    if err != nil {
        log.Fatal("Database connection failed: ", err)
    }

    // Create a new Gin router
    r := gin.Default()

    // Initialize repositories
    bookRepo := repository.NewBookRepository(db)
    userRepo := repository.NewUserRepository(db)
    borrowingRepo := repository.NewBorrowingRepository(db)

    // Initialize auth
    jwtAuth := auth.NewJWTAuth()

    // Initialize use cases
    bookUsecase := usecase.NewBookUsecase(bookRepo)
    userUsecase := usecase.NewUserUsecase(userRepo, jwtAuth)
    borrowingUsecase := usecase.NewBorrowingUsecase(borrowingRepo,bookRepo)

    // Set up handlers
    delivery.NewBookHandler(r, bookUsecase, jwtAuth)
    delivery.NewUserHandler(r, userUsecase)
    delivery.NewBorrowingHandler(r, borrowingUsecase,jwtAuth)

    r.Run()
}
