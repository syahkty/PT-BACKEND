package repository

import (
    "context"
    "library_api/internal/entity"

    "gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    return &user, err
}
