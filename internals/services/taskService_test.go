package services

import (
	"TaskSvc/commons/apploggers"
	"TaskSvc/internals/db"
	dbmodels "TaskSvc/internals/db/models"
	"TaskSvc/internals/models"
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertToModelTask(taskSchema *dbmodels.TaskSchema) *models.Task {
	return &models.Task{
		ID:          taskSchema.ID,
		Title:       taskSchema.Title,
		Description: taskSchema.Description,
		Status:      taskSchema.Status,
	}
}

var id = primitive.NewObjectID()

var _ = Describe("TaskService", func() {

	Describe("GetTaskById", func() {
		It("valid", func() {
			mockDbService := db.MockDbService{
				FakeGetTaskById: func(ctx context.Context, taskId string) (*dbmodels.TaskSchema, error) {
					return &dbmodels.TaskSchema{ID: id, Title: "Task 1", Description: "Task 1 Description", Status: "Pending"}, nil
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task, err := service.GetTaskById(ctx, "1")

			Expect(err).NotTo(HaveOccurred())
			Expect(task.Title).To(Equal("Task 1"))
		})

		It("error fetching task", func() {
			mockDbService := db.MockDbService{
				FakeGetTaskById: func(ctx context.Context, taskId string) (*dbmodels.TaskSchema, error) {
					return nil, fmt.Errorf("database error")
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task, err := service.GetTaskById(ctx, "1")

			Expect(err).To(HaveOccurred())
			Expect(task).To(BeNil())
		})
	})

	Describe("GetTasks", func() {
		It("valid", func() {
			mockDbService := db.MockDbService{
				FakeGetTasks: func(ctx context.Context) ([]*dbmodels.TaskSchema, error) {
					return []*dbmodels.TaskSchema{
						{ID: id, Title: "Task 1", Description: "Task 1 Description", Status: "Pending"},
						{ID: id, Title: "Task 2", Description: "Task 2 Description", Status: "Completed"},
					}, nil
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			tasks, err := service.GetTasks(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(tasks)).To(Equal(2))
			Expect(tasks[0].Title).To(Equal("Task 1"))
			Expect(tasks[1].Status).To(Equal("Completed"))
		})

		It("error fetching tasks", func() {
			mockDbService := db.MockDbService{
				FakeGetTasks: func(ctx context.Context) ([]*dbmodels.TaskSchema, error) {
					return nil, fmt.Errorf("database error")
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			tasks, err := service.GetTasks(ctx)

			Expect(err).To(HaveOccurred())
			Expect(tasks).To(BeNil())
		})
	})

	Describe("CreateTask", func() {
		It("valid", func() {
			mockDbService := db.MockDbService{
				FakeSaveTask: func(ctx context.Context, task *dbmodels.TaskSchema) (string, error) {
					return "1", nil
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task := &dbmodels.TaskSchema{Title: "New Task", Description: "New Task Description", Status: "Pending"}
			modelTask := convertToModelTask(task)

			taskId, err := service.CreateTask(ctx, modelTask)

			Expect(err).NotTo(HaveOccurred())
			Expect(taskId).To(Equal("1"))
		})

		It("error creating task", func() {
			mockDbService := db.MockDbService{
				FakeSaveTask: func(ctx context.Context, task *dbmodels.TaskSchema) (string, error) {
					return "", fmt.Errorf("database error")
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task := &dbmodels.TaskSchema{Title: "New Task", Description: "New Task Description", Status: "Pending"}
			modelTask := convertToModelTask(task)

			taskId, err := service.CreateTask(ctx, modelTask)

			Expect(err).To(HaveOccurred())
			Expect(taskId).To(Equal(""))
		})
	})

	Describe("UpdateTask", func() {
		It("valid", func() {
			mockDbService := db.MockDbService{
				FakeUpdateTask: func(ctx context.Context, task *dbmodels.TaskSchema, taskId string) error {
					return nil
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task := &dbmodels.TaskSchema{Title: "Updated Task", Description: "Updated Task Description", Status: "In-progress"}
			modelTask := convertToModelTask(task)

			err := service.UpdateTask(ctx, modelTask, "1")

			Expect(err).NotTo(HaveOccurred())
		})

		It("error updating task", func() {
			mockDbService := db.MockDbService{
				FakeUpdateTask: func(ctx context.Context, task *dbmodels.TaskSchema, taskId string) error {
					return fmt.Errorf("database error")
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			task := &dbmodels.TaskSchema{Title: "Updated Task", Description: "Updated Task Description", Status: "In-progress"}
			modelTask := convertToModelTask(task)

			err := service.UpdateTask(ctx, modelTask, "1")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteTaskById", func() {
		It("valid", func() {
			mockDbService := db.MockDbService{
				FakeDeleteTaskById: func(ctx context.Context, taskId string) error {
					return nil
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			err := service.DeleteTaskById(ctx, "1")

			Expect(err).NotTo(HaveOccurred())
		})

		It("error deleting task", func() {
			mockDbService := db.MockDbService{
				FakeDeleteTaskById: func(ctx context.Context, taskId string) error {
					return fmt.Errorf("database error")
				},
			}

			service := NewTaskService(mockDbService)
			ctx, _ := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

			err := service.DeleteTaskById(ctx, "1")

			Expect(err).To(HaveOccurred())
		})
	})
})
