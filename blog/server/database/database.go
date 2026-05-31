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
	initDefaultData()
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

func initDefaultData() {
	// 初始化默认分类
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "技术"},
			{Name: "生活"},
			{Name: "随笔"},
			{Name: "教程"},
		}
		DB.Create(&categories)
		log.Println("Default categories created")
	}

	// 初始化默认标签
	var tagCount int64
	DB.Model(&models.Tag{}).Count(&tagCount)
	if tagCount == 0 {
		tags := []models.Tag{
			{Name: "Go"},
			{Name: "Vue"},
			{Name: "JavaScript"},
			{Name: "编程"},
			{Name: "学习"},
			{Name: "分享"},
		}
		DB.Create(&tags)
		log.Println("Default tags created")
	}
}