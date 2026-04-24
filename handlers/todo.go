package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"todo-api/database"
	"todo-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	todo := models.Todo{
		Title:  req.Title,
		Status: 0,
	}

	if req.Description != "" {
		todo.Description = sql.NullString{String: req.Description, Valid: true}
	}

	if req.Status != nil {
		todo.Status = *req.Status
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			todo.DueDate = sql.NullTime{Time: dueDate, Valid: true}
		}
	}

	if err := database.GetDB().Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "创建失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "创建成功",
		Data:    todo,
	})
}

func GetTodos(c *gin.Context) {
	var query models.QueryTodoRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	db := database.GetDB()
	var todos []models.Todo

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Order("created_at DESC").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    todos,
	})
}

func GetTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "无效的ID",
			Data:    nil,
		})
		return
	}

	var todo models.Todo
	if err := database.GetDB().First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "记录不存在",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    todo,
	})
}

func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "无效的ID",
			Data:    nil,
		})
		return
	}

	var todo models.Todo
	if err := database.GetDB().First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "记录不存在",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			updates["due_date"] = dueDate
		}
	}

	if len(updates) > 0 {
		if err := database.GetDB().Model(&todo).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "更新失败: " + err.Error(),
				Data:    nil,
			})
			return
		}
	}

	database.GetDB().First(&todo, id)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "更新成功",
		Data:    todo,
	})
}

func DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "无效的ID",
			Data:    nil,
		})
		return
	}

	result := database.GetDB().Delete(&models.Todo{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "删除失败: " + result.Error.Error(),
			Data:    nil,
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "记录不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "删除成功",
		Data:    nil,
	})
}
