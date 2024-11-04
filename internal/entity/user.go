package entity

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username"`
    Password string `json:"password"`
    Role     string `json:"role"` //bisa "admin" atau "user"
}
