package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:50;uniqueIndex"`
	Password  string         `json:"-" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Articles  []Article      `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Articles  []Article      `json:"articles,omitempty" gorm:"many2many:article_tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Article struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Content      string         `json:"content" gorm:"type:text"`
	Status       int            `json:"status" gorm:"default:0"` // 0: draft, 1: published
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	IsTop        bool           `json:"is_top" gorm:"default:false"`
	Password     string         `json:"-" gorm:"size:255"`
	PublishedAt  *time.Time     `json:"published_at"`
	CategoryID   *uint          `json:"category_id"`
	Category     *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags         []Tag          `json:"tags,omitempty" gorm:"many2many:article_tags"`
	Comments     []Comment      `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Comment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ArticleID   uint           `json:"article_id" gorm:"not null;index"`
	Article     *Article       `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	AuthorName  string         `json:"author_name" gorm:"size:100;not null"`
	AuthorEmail string         `json:"author_email" gorm:"size:255;not null"`
	AuthorIP    string         `json:"author_ip" gorm:"size:50"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Likes       int            `json:"likes" gorm:"default:0"`
	IsAdmin     bool           `json:"is_admin" gorm:"default:false"`
	Status      int            `json:"status" gorm:"default:0"` // 0: pending, 1: approved, 2: deleted
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CommentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CommentID uint      `json:"comment_id" gorm:"not null;index"`
	IP        string    `json:"ip" gorm:"size:50;not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemSettings struct {
	ID                uint `json:"id" gorm:"primaryKey"`
	CommentModeration bool `json:"comment_moderation" gorm:"default:true"`
	CaptchaEnabled    bool `json:"captcha_enabled" gorm:"default:true"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	Status      int    `json:"status"`
	IsTop       bool   `json:"is_top"`
	Password    string `json:"password"`
	CategoryID  *uint  `json:"category_id"`
	TagIDs      []uint `json:"tag_ids"`
	PublishedAt string `json:"published_at"`
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

type CommentRequest struct {
	ArticleID   uint   `json:"article_id" binding:"required"`
	AuthorName  string `json:"author_name" binding:"required"`
	AuthorEmail string `json:"author_email" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Captcha     string `json:"captcha"`
	CaptchaID   string `json:"captcha_id"`
}

type CommentStatusRequest struct {
	Status int `json:"status" binding:"required"`
}

type SystemSettingsRequest struct {
	CommentModeration bool `json:"comment_moderation"`
	CaptchaEnabled    bool `json:"captcha_enabled"`
}
