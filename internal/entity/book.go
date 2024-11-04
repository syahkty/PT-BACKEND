package entity

type Book struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Title       string `json:"title"`
    Author      string `json:"author"`
    Description string `json:"description"`
    Year        int    `json:"year"`
}
