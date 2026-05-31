package controllers

import (
	"blog-server/database"
	"blog-server/models"
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var captchaStore = make(map[string]string)

func init() {
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			for k := range captchaStore {
				delete(captchaStore, k)
			}
		}
	}()
}

func GetArticleComments(c *gin.Context) {
	articleID := c.Param("id")

	var comments []models.Comment
	result := database.DB.Where("article_id = ? AND status = 1", articleID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    []models.Comment{},
		})
		return
	}

	for i := range comments {
		comments[i].AuthorEmail = getGravatarURL(comments[i].AuthorEmail)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    comments,
	})
}

func CreateComment(c *gin.Context) {
	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var settings models.SystemSettings
	database.DB.First(&settings)

	if settings.CaptchaEnabled {
		if req.Captcha == "" || req.CaptchaID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请输入验证码",
			})
			return
		}

		storedCaptcha, exists := captchaStore[req.CaptchaID]
		if !exists || storedCaptcha != req.Captcha {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "验证码错误或已过期",
			})
			return
		}
		delete(captchaStore, req.CaptchaID)
	}

	var article models.Article
	if err := database.DB.First(&article, req.ArticleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}

	isAdmin := false
	if username, exists := c.Get("username"); exists {
		if username == "admin" {
			isAdmin = true
		}
	}

	comment := models.Comment{
		ArticleID:   req.ArticleID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		AuthorIP:    c.ClientIP(),
		Content:     req.Content,
		IsAdmin:     isAdmin,
		Status:      1,
	}

	if !isAdmin && settings.CommentModeration {
		comment.Status = 0
	}

	database.DB.Create(&comment)

	if comment.Status == 1 {
		database.DB.Model(&article).Update("comment_count", article.CommentCount+1)
	}

	message := "评论发表成功"
	if comment.Status == 0 {
		message = "评论已提交，等待审核"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    comment,
	})
}

func LikeComment(c *gin.Context) {
	id := c.Param("id")
	ip := c.ClientIP()

	var existingLike models.CommentLike
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	result := database.DB.Where("comment_id = ? AND ip = ? AND created_at > ?", id, ip, oneDayAgo).First(&existingLike)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "您已在24小时内点赞过此评论",
		})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "评论不存在",
		})
		return
	}

	commentLike := models.CommentLike{
		CommentID: comment.ID,
		IP:        ip,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&commentLike)

	comment.Likes++
	database.DB.Save(&comment)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "点赞成功",
		"data": gin.H{
			"likes": comment.Likes,
		},
	})
}

func GetAdminComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.DefaultQuery("status", "")
	articleID := c.DefaultQuery("article_id", "")
	offset := (page - 1) * pageSize

	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if articleID != "" {
		query = query.Where("article_id = ?", articleID)
	}

	query.Count(&total)
	database.DB.Preload("Article").Where(query).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	for i := range comments {
		comments[i].AuthorEmail = getGravatarURL(comments[i].AuthorEmail)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":       comments,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func UpdateCommentStatus(c *gin.Context) {
	id := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "评论不存在",
		})
		return
	}

	var req models.CommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	oldStatus := comment.Status
	comment.Status = req.Status
	database.DB.Save(&comment)

	if oldStatus != 1 && req.Status == 1 {
		var article models.Article
		database.DB.First(&article, comment.ArticleID)
		database.DB.Model(&article).Update("comment_count", article.CommentCount+1)
	} else if oldStatus == 1 && req.Status != 1 {
		var article models.Article
		database.DB.First(&article, comment.ArticleID)
		if article.CommentCount > 0 {
			database.DB.Model(&article).Update("comment_count", article.CommentCount-1)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "状态更新成功",
		"data":    comment,
	})
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "评论不存在",
		})
		return
	}

	if comment.Status == 1 {
		var article models.Article
		database.DB.First(&article, comment.ArticleID)
		if article.CommentCount > 0 {
			database.DB.Model(&article).Update("comment_count", article.CommentCount-1)
		}
	}

	database.DB.Where("comment_id = ?", id).Delete(&models.CommentLike{})
	database.DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func DeleteCommentsBatch(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	for _, id := range req.IDs {
		var comment models.Comment
		if err := database.DB.First(&comment, id).Error; err == nil {
			if comment.Status == 1 {
				var article models.Article
				database.DB.First(&article, comment.ArticleID)
				if article.CommentCount > 0 {
					database.DB.Model(&article).Update("comment_count", article.CommentCount-1)
				}
			}
			database.DB.Where("comment_id = ?", id).Delete(&models.CommentLike{})
			database.DB.Delete(&comment)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量删除成功",
	})
}

func GetSystemSettings(c *gin.Context) {
	var settings models.SystemSettings
	database.DB.First(&settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    settings,
	})
}

func UpdateSystemSettings(c *gin.Context) {
	var req models.SystemSettingsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var settings models.SystemSettings
	database.DB.First(&settings)

	settings.CommentModeration = req.CommentModeration
	settings.CaptchaEnabled = req.CaptchaEnabled
	database.DB.Save(&settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设置更新成功",
		"data":    settings,
	})
}

func GenerateCaptcha(c *gin.Context) {
	num1 := rand.Intn(10) + 1
	num2 := rand.Intn(10) + 1
	result := fmt.Sprintf("%d", num1+num2)
	captchaID := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
	captchaStore[captchaID] = result

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"captcha_id": captchaID,
			"question":   fmt.Sprintf("%d + %d = ?", num1, num2),
		},
	})
}

func getGravatarURL(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	hash := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=48&d=identicon", hash)
}
