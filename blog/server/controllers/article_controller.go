package controllers

import (
	"blog-server/database"
	"blog-server/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetPublishedArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	database.DB.Model(&models.Article{}).Where("status = 1").Count(&total)
	database.DB.Preload("Category").Preload("Tags").Where("status = 1").Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":       articles,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	result := database.DB.Preload("Category").Preload("Tags").First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}

	article.ViewCount++
	database.DB.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    article,
	})
}

func GetAdminArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	database.DB.Model(&models.Article{}).Count(&total)
	database.DB.Preload("Category").Preload("Tags").Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":       articles,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}

func CreateArticle(c *gin.Context) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	article := models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Status:     req.Status,
		IsTop:      req.IsTop,
		Password:   req.Password,
		CategoryID: req.CategoryID,
	}

	if req.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.PublishedAt); err == nil {
			article.PublishedAt = &t
		}
	}

	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		database.DB.Find(&tags, req.TagIDs)
		article.Tags = tags
	}

	database.DB.Create(&article)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    article,
	})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	result := database.DB.First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Status = req.Status
	article.IsTop = req.IsTop
	article.Password = req.Password
	article.CategoryID = req.CategoryID

	if req.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.PublishedAt); err == nil {
			article.PublishedAt = &t
		}
	}

	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		database.DB.Find(&tags, req.TagIDs)
		database.DB.Model(&article).Association("Tags").Replace(tags)
	} else {
		database.DB.Model(&article).Association("Tags").Clear()
	}

	database.DB.Save(&article)
	database.DB.Preload("Category").Preload("Tags").First(&article, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    article,
	})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	result := database.DB.First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}

	database.DB.Model(&article).Association("Tags").Clear()
	database.DB.Delete(&article)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择文件",
		})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件上传失败",
		})
		return
	}

	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": url,
		},
	})
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    categories,
	})
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	database.DB.Create(&category)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Category{}, id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func GetTags(c *gin.Context) {
	var tags []models.Tag
	database.DB.Find(&tags)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    tags,
	})
}

func CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	database.DB.Create(&tag)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    tag,
	})
}

func DeleteTag(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Tag{}, id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
