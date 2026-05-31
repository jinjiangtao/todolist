package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:50;uniqueIndex"`
	Password  string         `json:"-" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Article struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text"`
	Status    int            `json:"status" gorm:"default:0"`
	ViewCount int            `json:"view_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
