package usecase

import (
	"context"
	"library_api/internal/entity"
	"library_api/internal/repository"
	"library_api/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
    userRepo repository.UserRepository
	auth     *auth.JWTAuth
}

func NewUserUsecase(repo repository.UserRepository, auth *auth.JWTAuth) *UserUsecase {
    return &UserUsecase{userRepo: repo, auth: auth}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user *entity.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return u.userRepo.Create(ctx, user)
}

func (u *UserUsecase) LoginUser(ctx context.Context, username, password string) (string, error) {
    user, err := u.userRepo.GetByUsername(ctx, username)
    if err != nil {
        return "", err
    }

    // jika password matches
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", err
    }

    // Generate JWT token using JWTAuth method
    token, err := u.auth.GenerateJWT(user.Username, user.Role,user.ID)
    if err != nil {
        return "", err
    }
    return token, nil
}
