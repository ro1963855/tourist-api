package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(POSTGRES_HOST, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT string) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"password:email"`
	UpdatedAt string `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt string `gorm:"column:created_at;autoCreateTime"`
}

func AuthenticateUser(email string) (User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
