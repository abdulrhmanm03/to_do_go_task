package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"to_do_api/controllers"
	"to_do_api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestTaskDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func newTestTaskRouter(userID string) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	return r
}

func TestCreateTask(t *testing.T) {
	db := setupTestTaskDB(t)
	userID := uuid.New().String()
	router := newTestTaskRouter(userID)
	router.POST("/tasks", controllers.CreateTask(db))

	taskBody := map[string]interface{}{
		"title":       "Test Task",
		"description": "Test Description",
		"completed":   false,
	}
	bodyBytes, err := json.Marshal(taskBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	db := setupTestTaskDB(t)
	userID := uuid.New().String()
	router := newTestTaskRouter(userID)
	router.POST("/tasks", controllers.CreateTask(db))

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListTasks(t *testing.T) {
	db := setupTestTaskDB(t)

	userID := uuid.New()
	otherUserID := uuid.New()
	tasks := []models.Task{
		{Title: "Task 1", UserID: userID},
		{Title: "Task 2", UserID: userID},
		{Title: "Other Task", UserID: otherUserID},
	}
	for _, task := range tasks {
		err := db.Create(&task).Error
		assert.NoError(t, err)
	}

	router := newTestTaskRouter(userID.String())
	router.GET("/tasks", controllers.ListTasks(db))

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respTasks []models.Task
	err = json.Unmarshal(w.Body.Bytes(), &respTasks)
	assert.NoError(t, err)
	assert.Len(t, respTasks, 2)
}

func TestUpdateTask(t *testing.T) {
	db := setupTestTaskDB(t)
	userID := uuid.New()

	task := models.Task{
		Title:  "Original Title",
		UserID: userID,
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	router := newTestTaskRouter(userID.String())
	router.PUT("/tasks/:id", controllers.UpdateTask(db))

	updateBody := map[string]interface{}{
		"title": "Updated Title",
	}
	bodyBytes, err := json.Marshal(updateBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/tasks/"+task.ID.String(), bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateTask_Unauthorized(t *testing.T) {
	db := setupTestTaskDB(t)
	// Create a task that belongs to another user.
	userID := uuid.New()
	otherUserID := uuid.New()
	task := models.Task{
		Title:  "Original Title",
		UserID: otherUserID,
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	// Set request context with a different user.
	router := newTestTaskRouter(userID.String())
	router.PUT("/tasks/:id", controllers.UpdateTask(db))

	updateBody := map[string]interface{}{
		"title": "Updated Title",
	}
	bodyBytes, err := json.Marshal(updateBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/tasks/"+task.ID.String(), bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateTask_InvalidID(t *testing.T) {
	db := setupTestTaskDB(t)
	router := newTestTaskRouter(uuid.New().String())
	router.PUT("/tasks/:id", controllers.UpdateTask(db))

	updateBody := map[string]interface{}{
		"title": "Updated Title",
	}
	bodyBytes, err := json.Marshal(updateBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/tasks/invalid-uuid", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTask(t *testing.T) {
	db := setupTestTaskDB(t)
	userID := uuid.New()

	task := models.Task{
		Title:  "Task to delete",
		UserID: userID,
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	router := newTestTaskRouter(userID.String())
	router.DELETE("/tasks/:id", controllers.DeleteTask(db))

	req, err := http.NewRequest("DELETE", "/tasks/"+task.ID.String(), nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Task deleted successfully", resp["message"])

	var deletedTask models.Task
	err = db.First(&deletedTask, task.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteTask_Unauthorized(t *testing.T) {
	db := setupTestTaskDB(t)
	// Create a task that belongs to a different user.
	userID := uuid.New()
	otherUserID := uuid.New()
	task := models.Task{
		Title:  "Task not allowed to delete",
		UserID: otherUserID,
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	router := newTestTaskRouter(userID.String())
	router.DELETE("/tasks/:id", controllers.DeleteTask(db))

	req, err := http.NewRequest("DELETE", "/tasks/"+task.ID.String(), nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteTask_InvalidID(t *testing.T) {
	db := setupTestTaskDB(t)
	router := newTestTaskRouter(uuid.New().String())
	router.DELETE("/tasks/:id", controllers.DeleteTask(db))

	req, err := http.NewRequest("DELETE", "/tasks/invalid-uuid", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
