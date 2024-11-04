package repository

import (
    "context"
    "library_api/internal/entity"
    "gorm.io/gorm"
)

type BorrowingRepository interface {
    GetBorrowedBooksByUserID(ctx context.Context, userID uint) ([]entity.Borrowing, error)
    GetAllBorrowings(ctx context.Context) ([]entity.Borrowing, error) // Method untuk admin
    Create(ctx context.Context, borrowing *entity.Borrowing) error    // Method untuk membuat record peminjaman baru
    GetCountByUserAndBook(ctx context.Context, userID uint, bookID uint, count *int64) error // Method untuk mengecek peminjaman buku aktif oleh user tertentu
    GetByID(ctx context.Context, borrowingID uint) (*entity.Borrowing, error)             // Method untuk mendapatkan peminjaman berdasarkan ID
    Update(ctx context.Context, borrowing *entity.Borrowing) error   
    Delete(ctx context.Context, borrowingID uint) error // Method untuk menghapus peminjaman
    GetByBookAndUserID(ctx context.Context, bookID uint, userID uint) (*entity.Borrowing, error) // Method untuk mendapatkan peminjaman berdasarkan ID buku dan ID pengguna
}

type borrowingRepository struct {
    db *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) BorrowingRepository {
    return &borrowingRepository{db}
}

func (r *borrowingRepository) GetBorrowedBooksByUserID(ctx context.Context, userID uint) ([]entity.Borrowing, error) {
    var borrowings []entity.Borrowing
    err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&borrowings).Error
    return borrowings, err
}

func (r *borrowingRepository) GetAllBorrowings(ctx context.Context) ([]entity.Borrowing, error) {
    var borrowings []entity.Borrowing
    err := r.db.WithContext(ctx).Find(&borrowings).Error
    return borrowings, err
}

func (r *borrowingRepository) GetCountByUserAndBook(ctx context.Context, userID uint, bookID uint, count *int64) error {
    return r.db.WithContext(ctx).
        Model(&entity.Borrowing{}).
        Where("user_id = ? AND book_id = ? AND status = ?", userID, bookID, "active").
        Count(count).Error
}

func (r *borrowingRepository) Create(ctx context.Context, borrowing *entity.Borrowing) error {
    return r.db.WithContext(ctx).Create(borrowing).Error
}

func (r *borrowingRepository) GetByID(ctx context.Context, borrowingID uint) (*entity.Borrowing, error) {
    var borrowing entity.Borrowing
    if err := r.db.WithContext(ctx).First(&borrowing, borrowingID).Error; err != nil {
        return nil, err
    }
    return &borrowing, nil
}

func (r *borrowingRepository) Update(ctx context.Context, borrowing *entity.Borrowing) error {
    return r.db.WithContext(ctx).Save(borrowing).Error
}

func (r *borrowingRepository) Delete(ctx context.Context, borrowingID uint) error {
    return r.db.WithContext(ctx).Delete(&entity.Borrowing{}, borrowingID).Error // Menghapus record peminjaman berdasarkan ID
}

func (r *borrowingRepository) GetByBookAndUserID(ctx context.Context, bookID uint, userID uint) (*entity.Borrowing, error) {
    var borrowing entity.Borrowing
    err := r.db.WithContext(ctx).
        Where("book_id = ? AND user_id = ? AND status = ?", bookID, userID, "active").
        First(&borrowing).Error
    if err != nil {
        return nil, err
    }
    return &borrowing, nil
}
