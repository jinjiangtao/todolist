package database

import (
	"blog-server/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"

	sqlite "github.com/glebarez/sqlite"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Category{}, &models.Tag{})

	initDefaultAdmin()
}

func initDefaultAdmin() {
	var user models.User
	result := DB.Where("username = ?", "admin").First(&user)
	if result.Error != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		admin := models.User{Username: "admin", Password: string(hashedPassword)}
		DB.Create(&admin)
		log.Println("Default admin account created: admin/123456")
		log.Printf("Password hash: %s (length: %d)", string(hashedPassword), len(string(hashedPassword)))
	} else {
		log.Println("Admin account already exists")
		log.Printf("Existing password hash length: %d", len(user.Password))
	}
}