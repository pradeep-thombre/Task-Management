package commons

import (
	dbmodels "TaskSvc/internals/db/models"
	"TaskSvc/internals/models"
)

func MapToModel(taskSchema *dbmodels.TaskSchema) *models.Task {
	return &models.Task{
		ID:          taskSchema.ID,
		Title:       taskSchema.Title,
		Description: taskSchema.Description,
		Status:      taskSchema.Status,
		CreatedAt:   taskSchema.CreatedAt,
		UpdatedAt:   taskSchema.UpdatedAt,
	}
}
