package db

import (
	dbmodels "TaskSvc/internals/db/models"
	"context"
	"fmt"
)

type MockDbService struct {
	FakeGetTaskById    func(ctx context.Context, taskId string) (*dbmodels.TaskSchema, error)
	FakeSaveTask       func(ctx context.Context, task *dbmodels.TaskSchema) (string, error)
	FakeUpdateTask     func(ctx context.Context, task *dbmodels.TaskSchema, taskId string) error
	FakeDeleteTaskById func(ctx context.Context, taskId string) error
	FakeGetTasks       func(ctx context.Context) ([]*dbmodels.TaskSchema, error)
}

func (m MockDbService) GetTaskById(ctx context.Context, taskId string) (*dbmodels.TaskSchema, error) {
	if m.FakeGetTaskById != nil {
		return m.FakeGetTaskById(ctx, taskId)
	}
	return nil, fmt.Errorf("GetTaskById-error")
}

func (m MockDbService) SaveTask(ctx context.Context, task *dbmodels.TaskSchema) (string, error) {
	if m.FakeSaveTask != nil {
		return m.FakeSaveTask(ctx, task)
	}
	return "", fmt.Errorf("SaveTask-error")
}

func (m MockDbService) UpdateTask(ctx context.Context, task *dbmodels.TaskSchema, taskId string) error {
	if m.FakeUpdateTask != nil {
		return m.FakeUpdateTask(ctx, task, taskId)
	}
	return fmt.Errorf("UpdateTask-error")
}

func (m MockDbService) DeleteTaskById(ctx context.Context, taskId string) error {
	if m.FakeDeleteTaskById != nil {
		return m.FakeDeleteTaskById(ctx, taskId)
	}
	return fmt.Errorf("DeleteTaskById-error")
}

func (m MockDbService) GetTasks(ctx context.Context) ([]*dbmodels.TaskSchema, error) {
	if m.FakeGetTasks != nil {
		return m.FakeGetTasks(ctx)
	}
	return nil, fmt.Errorf("GetTasks-error")
}
