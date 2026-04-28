package models

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description sql.NullString `gorm:"type:text" json:"description"`
	Status      int            `gorm:"type:tinyint;not null;default:0" json:"status"`
	DueDate     sql.NullTime   `gorm:"type:date" json:"due_date"`
	UID         int64          `gorm:"type:int;not null" json:"uid"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Todo) TableName() string {
	return "todos"
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
	DueDate     string `json:"due_date"`
	UID         int64  `json:"uid" binding:"required"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
	DueDate     string `json:"due_date"`
	UID         int64  `json:"uid"`
}

type QueryTodoRequest struct {
	Status   *int `form:"status"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=10"`
}

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(255);not null;unique" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Bio       string    `gorm:"type:text" json:"bio"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Bio      string `json:"bio"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
