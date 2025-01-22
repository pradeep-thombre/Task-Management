package services

import (
	"TaskSvc/commons"
	"TaskSvc/commons/apploggers"
	"TaskSvc/internals/db"
	dbmodel "TaskSvc/internals/db/models"
	"TaskSvc/internals/models"
	"context"
	"encoding/json"
)

type TaskService interface {
	GetTaskById(context context.Context, taskId string) (*models.Task, error)
	DeleteTaskById(context context.Context, taskId string) error
	GetTasks(context context.Context) ([]*models.Task, error)
	CreateTask(context context.Context, task *models.Task) (string, error)
	UpdateTask(context context.Context, task *models.Task, taskId string) error
}

type taskService struct {
	dbservice db.DbService
}

func NewTaskService(dbservice db.DbService) TaskService {
	return &taskService{dbservice: dbservice}
}

func (s *taskService) GetTaskById(ctx context.Context, taskId string) (*models.Task, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	taskSchema, err := s.dbservice.GetTaskById(ctx, taskId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return commons.MapToModel(taskSchema), nil
}

func (s *taskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	taskSchemas, err := s.dbservice.GetTasks(ctx)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tasks := make([]*models.Task, len(taskSchemas))
	for i, taskSchema := range taskSchemas {
		tasks[i] = commons.MapToModel(taskSchema)
	}

	return tasks, nil
}

func (s *taskService) DeleteTaskById(ctx context.Context, taskId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	if err := s.dbservice.DeleteTaskById(ctx, taskId); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *taskService) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	var taskSchema dbmodel.TaskSchema
	pbyes, _ := json.Marshal(task)
	err := json.Unmarshal(pbyes, &taskSchema)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	taskId, err := s.dbservice.SaveTask(ctx, &taskSchema)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return taskId, nil
}

func (s *taskService) UpdateTask(ctx context.Context, task *models.Task, taskId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	var taskSchema dbmodel.TaskSchema
	pbyes, _ := json.Marshal(task)
	err := json.Unmarshal(pbyes, &taskSchema)
	if err != nil {
		logger.Error(err)
		return err
	}
	return s.dbservice.UpdateTask(ctx, &taskSchema, taskId)
}
