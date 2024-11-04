package configs

import (
	"fmt"
	"library_api/internal/entity"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
    // Memuat file .env
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("gagal memuat file .env: %v", err)
    }

    // Membuat DSN dari variabel lingkungan
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_CHARSET"),
        os.Getenv("DB_PARSE_TIME"),
        os.Getenv("DB_LOC"),
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // AutoMigrate untuk membuat tabel secara otomatis
    db.AutoMigrate(&entity.User{}, &entity.Book{}, &entity.Borrowing{})

    return db, nil
}