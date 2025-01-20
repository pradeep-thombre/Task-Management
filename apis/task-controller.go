package apis

import (
	"TaskSvc/commons"
	"TaskSvc/internals/models"
	"TaskSvc/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (t *TaskController) GetTasks(c *gin.Context) {
	tasks, err := t.taskService.GetTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to fetch tasks", nil))
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) GetTaskById(c *gin.Context) {
	taskId := c.Param("id")
	task, err := t.taskService.GetTaskById(c, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to fetch task", nil))
		return
	}
	c.JSON(http.StatusOK, task)
}

func (t *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request payload", nil))
		return
	}

	taskId, err := t.taskService.CreateTask(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to create task", nil))
		return
	}
	c.JSON(http.StatusCreated, map[string]string{"id": taskId})
}

func (t *TaskController) UpdateTask(c *gin.Context) {
	taskId := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request payload", nil))
		return
	}

	if err := t.taskService.UpdateTask(c, &task, taskId); err != nil {
		c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to update task", nil))
		return
	}
	c.Status(http.StatusOK)
}

func (t *TaskController) DeleteTask(c *gin.Context) {
	taskId := c.Param("id")
	if err := t.taskService.DeleteTaskById(c, taskId); err != nil {
		c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to delete task", nil))
		return
	}
	c.Status(http.StatusNoContent)
}
