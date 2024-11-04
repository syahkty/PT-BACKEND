package repository

import (
    "context"
    "library_api/internal/entity"

    "gorm.io/gorm"
)

type BookRepository interface {
    FindBooks(limit int, offset int, title string) ([]entity.Book, int64, error)
    Create(ctx context.Context, book *entity.Book) error
    GetAll(ctx context.Context) ([]entity.Book, error)
    GetByID(ctx context.Context, id uint) (*entity.Book, error)
    Update(ctx context.Context, book *entity.Book) error
    Delete(ctx context.Context, id uint) error
}

type bookRepository struct {
    db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
    return &bookRepository{db}
}

func (r *bookRepository) FindBooks(limit int, offset int, title string) ([]entity.Book, int64, error) {
	var books []entity.Book
	var total int64

	// Query dasar
	query := r.db.Model(&entity.Book{})

	// Filtering berdasarkan judul buku (opsional)
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// Hitung total hasil sebelum pagination
	query.Count(&total)

	// Pagination
	err := query.Limit(limit).Offset(offset).Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
    return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) GetAll(ctx context.Context) ([]entity.Book, error) {
    var books []entity.Book
    err := r.db.WithContext(ctx).Find(&books).Error
    return books, err
}

func (r *bookRepository) GetByID(ctx context.Context, id uint) (*entity.Book, error) {
    var book entity.Book
    err := r.db.WithContext(ctx).First(&book, id).Error
    return &book, err
}

func (r *bookRepository) Update(ctx context.Context, book *entity.Book) error {
    var existingBook entity.Book
    if err := r.db.WithContext(ctx).First(&existingBook, book.ID).Error; err != nil {
        return err 
    }
    return r.db.WithContext(ctx).Model(&existingBook).Updates(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&entity.Book{}, id).Error
}
