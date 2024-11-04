package configs

import (
    "library_api/internal/entity"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
    dsn := "root:@tcp(127.0.0.1:3306)/library_db?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // AutoMigrate untuk membuat tabel secara otomatis
    db.AutoMigrate(&entity.User{}, &entity.Book{}, &entity.Borrowing{})

    return db, nil
}   
