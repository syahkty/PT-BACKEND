package entity

import "time"

type Borrowing struct {
    ID         uint      `json:"id" gorm:"primaryKey"`  
    BookID     uint      `json:"book_id"`
    BookTitle  string    `json:"book_title"`               
    UserID     uint      `json:"user_id"`               
    BorrowedAt time.Time `json:"borrowed_at"`           
    ReturnedAt *time.Time `json:"returned_at,omitempty"` 
    Status     string     `json:"status"`                
}
