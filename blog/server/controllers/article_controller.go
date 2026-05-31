package controllers

import (
	"blog-server/database"
	"blog-server/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

func GetPublishedArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	database.DB.Model(&models.Article{}).Where("status = 1").Count(&total)
	database.DB.Where("status = 1").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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

	result := database.DB.First(&article, id)
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
	database.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
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
	database.DB.Save(&article)

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

	database.DB.Delete(&article)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
