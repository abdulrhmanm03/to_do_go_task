package main

import (
	"log"
	"to_do_api/auth"
	"to_do_api/config"
	"to_do_api/controllers"
	"to_do_api/database"
	"to_do_api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	r := gin.Default()

	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db, &auth.DefaultAuthService{}))

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/tasks", controllers.CreateTask(db))
		authorized.GET("/tasks", controllers.ListTasks(db))
		authorized.PUT("/tasks/:id", controllers.UpdateTask(db))
		authorized.DELETE("/tasks/:id", controllers.DeleteTask(db))
	}

	log.Fatal(r.Run(":" + cfg.PORT))
}
