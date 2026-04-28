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
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
	DueDate     string `json:"due_date"`
}

type QueryTodoRequest struct {
	Status   *int `form:"status"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=10"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
