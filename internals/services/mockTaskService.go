package services

import (
	"TaskSvc/internals/models"
	"context"
	"fmt"
)

type MockTaskService struct {
	FakeGetTaskById    func(ctx context.Context, taskId string) (*models.Task, error)
	FakeGetTasks       func(ctx context.Context) ([]*models.Task, error)
	FakeCreateTask     func(ctx context.Context, task *models.Task) (string, error)
	FakeUpdateTask     func(ctx context.Context, task *models.Task, taskId string) error
	FakeDeleteTaskById func(ctx context.Context, taskId string) error
}

func (m MockTaskService) GetTaskById(ctx context.Context, taskId string) (*models.Task, error) {
	if m.FakeGetTaskById != nil {
		return m.FakeGetTaskById(ctx, taskId)
	}
	return nil, fmt.Errorf("GetTaskById-error")
}

func (m MockTaskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	if m.FakeGetTasks != nil {
		return m.FakeGetTasks(ctx)
	}
	return nil, fmt.Errorf("GetTasks-error")
}

func (m MockTaskService) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	if m.FakeCreateTask != nil {
		return m.FakeCreateTask(ctx, task)
	}
	return "", fmt.Errorf("CreateTask-error")
}

func (m MockTaskService) UpdateTask(ctx context.Context, task *models.Task, taskId string) error {
	if m.FakeUpdateTask != nil {
		return m.FakeUpdateTask(ctx, task, taskId)
	}
	return fmt.Errorf("UpdateTask-error")
}

func (m MockTaskService) DeleteTaskById(ctx context.Context, taskId string) error {
	if m.FakeDeleteTaskById != nil {
		return m.FakeDeleteTaskById(ctx, taskId)
	}
	return fmt.Errorf("DeleteTaskById-error")
}
