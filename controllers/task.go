package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"to_do_api/models"
)

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))
		task.UserID = userID

		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

func ListTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		userID, _ := uuid.Parse(c.GetString("user_id"))

		if err := db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))
		if task.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this task"})
			return
		}

		var updateData models.Task
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Model(&task).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))
		if task.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this task"})
			return
		}

		if err := db.Delete(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	}
}
