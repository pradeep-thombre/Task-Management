package main

import (
	"TaskSvc/apis"
	"TaskSvc/apis/middleware"
	"TaskSvc/commons/apploggers"
	"TaskSvc/configs"
	"TaskSvc/internals/db"
	"TaskSvc/internals/services"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	context, logger := apploggers.NewLoggerWithCorrelationid(context.Background(), "")
	err := configs.NewApplicationConfig(context)
	if err != nil {
		logger.Errorf("Error in Appconfig:", err)
	}

	dbService := db.NewDbService(configs.AppConfig.DbClient)
	taskService := services.NewTaskService(dbService)
	taskController := apis.NewTaskController(taskService)

	// Initialize Gin router
	r := gin.Default()

	r.GET("/public/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", middleware.AuthenticateJWT, taskController.GetTaskById)
	r.POST("/tasks", middleware.AuthenticateJWT, taskController.CreateTask)
	r.PUT("/tasks/:id", middleware.AuthenticateJWT, taskController.UpdateTask)
	r.DELETE("/tasks/:id", middleware.AuthenticateJWT, taskController.DeleteTask)

	// Run the server
	r.Run(":8080")
}
