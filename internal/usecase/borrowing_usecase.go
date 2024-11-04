package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"library_api/internal/entity"
	"library_api/internal/repository"
	"net/http"
	"time"
)

type BorrowingUsecase struct {
    borrowingRepo repository.BorrowingRepository
    bookRepo      repository.BookRepository // Tambahkan repository untuk buku
}

func NewBorrowingUsecase(borrowingRepo repository.BorrowingRepository, bookRepo repository.BookRepository) *BorrowingUsecase {
    return &BorrowingUsecase{
        borrowingRepo: borrowingRepo,
        bookRepo:      bookRepo,
    }
}

func (u *BorrowingUsecase) GetBorrowedBooksByUserID(ctx context.Context, userID uint) ([]entity.Borrowing, error) {
    return u.borrowingRepo.GetBorrowedBooksByUserID(ctx, userID)
}

func (u *BorrowingUsecase) GetAllBorrowings(ctx context.Context) ([]entity.Borrowing, error) {
    return u.borrowingRepo.GetAllBorrowings(ctx)
}

func (u *BorrowingUsecase) CreateBorrowing(ctx context.Context, borrowing *entity.Borrowing) error {
    // Cek apakah buku tersedia
    book, err := u.bookRepo.GetByID(ctx, borrowing.BookID)
    if err != nil {
        return errors.New("book not found")
    }

    var activeBorrowingCount int64
    err = u.borrowingRepo.GetCountByUserAndBook(ctx, borrowing.UserID, borrowing.BookID, &activeBorrowingCount)
    if err != nil {
        return err
    }

    if activeBorrowingCount > 0 {
        return errors.New("this book is already borrowed by the user and not returned yet")
    }
    bookTitle := book.Title
    // Set waktu peminjaman dan status
    borrowing.BookTitle = bookTitle
    borrowing.BorrowedAt = time.Now()
    borrowing.Status = "active"

    if err := u.bookRepo.Update(ctx, book); err != nil {
        return err
    }

    // Simpan record peminjaman
    return u.borrowingRepo.Create(ctx, borrowing)
}

func (u *BorrowingUsecase) CreateBorrowingExternal(ctx context.Context, borrowing *entity.Borrowing) error {
    url := fmt.Sprintf("https://www.dbooks.org/api/book/%d", borrowing.BookID)
    fmt.Println("Checking availability for Book ID:", borrowing.BookID)
    fmt.Println("URL:", url)
    resp, err := http.Get(url)
    if err != nil {
        return errors.New("error checking book availability")
    }
    
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Membaca isi dari resp.Body untuk ditampilkan
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return errors.New("error reading response body")
        }
        return fmt.Errorf("book not available: %s", body)
    }

    var apiResponse struct {
        ID string `json:"id"`
        Title string `json:"title"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        return err
    }
    
    // Ambil ID buku
    bookID := apiResponse.ID
    bookTitle := apiResponse.Title
    fmt.Println("ID Buku:", bookID,bookTitle)

    var activeBorrowingCount int64
    err = u.borrowingRepo.GetCountByUserAndBook(ctx, borrowing.UserID, borrowing.BookID, &activeBorrowingCount)
    if err != nil {
        return err
    }

    if activeBorrowingCount > 0 {
        return errors.New("this book is already borrowed by the user and not returned yet")
    }

    // Set waktu peminjaman dan status
    borrowing.BookTitle = bookTitle
    borrowing.BorrowedAt = time.Now()
    borrowing.Status = "active"

    // Simpan record peminjaman
    return u.borrowingRepo.Create(ctx, borrowing)
}



func (u *BorrowingUsecase) ReturnBook(ctx context.Context, bookID uint, userID uint) error {
    // Cari record peminjaman berdasarkan BookID dan UserID
    borrowing, err := u.borrowingRepo.GetByBookAndUserID(ctx, bookID, userID)
    if err != nil {
        return errors.New("borrowing record not found")
    }

    if borrowing.UserID != userID {
        return errors.New("user ID does not match borrowing record")
    }
    // Cek apakah status sudah "returned"
    if borrowing.Status == "returned" {
        return errors.New("book is already returned")
    }

    // Set waktu pengembalian dan ubah status
    now := time.Now()
    borrowing.ReturnedAt = &now
    borrowing.Status = "returned"

    if err := u.borrowingRepo.Update(ctx, borrowing); err != nil {
        return err
    }

    return u.borrowingRepo.Delete(ctx, borrowing.ID)
}
