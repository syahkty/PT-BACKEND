package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"library_api/internal/entity"
	"library_api/internal/repository"
	"net/http"
	"net/url"
)

type BookUsecase struct {
    bookRepo repository.BookRepository
}

func NewBookUsecase(repo repository.BookRepository) *BookUsecase {
    return &BookUsecase{bookRepo: repo}
}

func (u *BookUsecase) CreateBook(ctx context.Context, book *entity.Book) error {
    return u.bookRepo.Create(ctx, book)
}

func (u *BookUsecase) GetBooks(limit int, page int, title string) ([]entity.Book, int64, error) {
	offset := (page - 1) * limit
	books, total, err := u.bookRepo.FindBooks(limit, offset, title)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (u *BookUsecase) GetBooksFromExternalAPI(ctx context.Context, title string) (interface{}, error) {
    apiURL := fmt.Sprintf("https://www.dbooks.org/api/search/%s", url.QueryEscape(title))
    resp, err := http.Get(apiURL)
if err != nil {
    return nil, err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    return nil, err
}

if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("error from API: %s", body)
}

// Decode JSON dari body
var result interface{}
if err := json.Unmarshal(body, &result); err != nil {
    return nil, err
}

return result, nil
}

func (u *BookUsecase) GetBookByID(ctx context.Context, id uint) (*entity.Book, error) {
    return u.bookRepo.GetByID(ctx, id)
}

func (u *BookUsecase) UpdateBook(ctx context.Context, book *entity.Book) error {
    return u.bookRepo.Update(ctx, book)
}

func (u *BookUsecase) DeleteBook(ctx context.Context, id uint) error {
    return u.bookRepo.Delete(ctx, id)
}
