package apis

import (
	"TaskSvc/commons"
	"TaskSvc/internals/models"
	"TaskSvc/internals/services"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task API Controller", func() {

	Describe("GetTasks", func() {
		It("valid", func() {
			eservice := services.MockTaskService{
				FakeGetTasks: func(ctx context.Context) ([]*models.Task, error) {
					return []*models.Task{
						{
							Title:       "Task 1",
							Description: "Task 1 Description",
							Status:      "Pending",
						},
					}, nil
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(eservice)
			controller.GetTasks(c)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response map[string]interface{} // Expect a map with "total" and "tasks" keys
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())

			// Check that the response contains the correct structure
			Expect(response).To(HaveKey("total"))
			Expect(response).To(HaveKey("tasks"))

			// Assert that "total" is correct
			total := response["total"].(float64) // JSON numbers are unmarshalled as float64
			Expect(total).To(Equal(1.0))

			// Assert that the "tasks" field is a slice and check its contents
			tasks, ok := response["tasks"].([]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(tasks)).To(Equal(1))

			task := tasks[0].(map[string]interface{})
			Expect(task["title"]).To(Equal("Task 1"))
			Expect(task["description"]).To(Equal("Task 1 Description"))
			Expect(task["status"]).To(Equal("Pending"))
		})

		It("error fetching tasks", func() {
			eservice := services.MockTaskService{
				FakeGetTasks: func(ctx context.Context) ([]*models.Task, error) {
					return nil, fmt.Errorf("failed to fetch tasks")
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(eservice)
			controller.GetTasks(c)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Failed to fetch tasks"))
		})
	})

	Describe("GetTaskById", func() {
		It("valid", func() {
			eservice := services.MockTaskService{
				FakeGetTaskById: func(ctx context.Context, taskId string) (*models.Task, error) {
					return &models.Task{
						ID:          taskId,
						Title:       "Task 1",
						Description: "Task 1 Description",
						Status:      "Pending",
					}, nil
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(eservice)
			controller.GetTaskById(c)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response *models.Task
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Title).To(Equal("Task 1"))
			Expect(response.Description).To(Equal("Task 1 Description"))
			Expect(response.Status).To(Equal("Pending"))
		})

		It("error fetching task", func() {
			eservice := services.MockTaskService{
				FakeGetTaskById: func(ctx context.Context, taskId string) (*models.Task, error) {
					return nil, fmt.Errorf("failed to fetch task")
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(eservice)
			controller.GetTaskById(c)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Failed to fetch task"))
		})

		It("missing task id", func() {
			eservice := services.MockTaskService{
				FakeGetTaskById: func(ctx context.Context, taskId string) (*models.Task, error) {
					return nil, nil
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks/", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: ""}} // Missing task id

			controller := NewTaskController(eservice)
			controller.GetTaskById(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Task ID is required"))
		})
	})

	Describe("CreateTask", func() {
		It("valid", func() {
			eservice := services.MockTaskService{
				FakeCreateTask: func(ctx context.Context, task *models.Task) (string, error) {
					return "test-id", nil
				},
			}
			task := &models.Task{
				Title:       "New Task",
				Description: "New Task Description",
				Status:      "Pending",
			}
			pbytes, err := json.Marshal(task)
			Expect(err).To(BeNil())
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(pbytes))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(eservice)
			controller.CreateTask(c)

			Expect(rec.Code).To(Equal(http.StatusCreated))

			var response map[string]string
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response["id"]).To(Equal("test-id"))
		})

		It("invalid payload", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("invalid"))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(services.MockTaskService{})
			controller.CreateTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Invalid request payload"))
		})

		It("missing Title", func() {
			task := &models.Task{
				Title:       "",
				Description: "New Task Description",
				Status:      "Pending",
			}
			pbytes, err := json.Marshal(task)
			Expect(err).To(BeNil())
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(pbytes))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(services.MockTaskService{})
			controller.CreateTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Title is required"))
		})

		It("missing Description", func() {
			task := &models.Task{
				Title:       "New Task",
				Description: "",
				Status:      "Pending",
			}
			pbytes, err := json.Marshal(task)
			Expect(err).To(BeNil())
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(pbytes))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			controller := NewTaskController(services.MockTaskService{})
			controller.CreateTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Description is required"))
		})
	})

	Describe("UpdateTask", func() {
		It("valid", func() {
			eservice := services.MockTaskService{
				FakeUpdateTask: func(ctx context.Context, task *models.Task, taskId string) error {
					return nil
				},
			}
			task := &models.Task{
				Title:       "Updated Task",
				Description: "Updated Task Description",
				Status:      "In-progress",
			}
			pbytes, err := json.Marshal(task)
			Expect(err).To(BeNil())
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(pbytes))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(eservice)
			controller.UpdateTask(c)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("invalid payload", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader("invalid"))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(services.MockTaskService{})
			controller.UpdateTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Invalid request payload"))
		})

		It("missing Title", func() {
			task := &models.Task{
				Title:       "",
				Description: "Updated Task Description",
				Status:      "In-progress",
			}
			pbytes, err := json.Marshal(task)
			Expect(err).To(BeNil())
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(pbytes))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(services.MockTaskService{})
			controller.UpdateTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Title is required"))
		})
	})

	Describe("DeleteTask", func() {
		It("valid", func() {
			eservice := services.MockTaskService{
				FakeDeleteTaskById: func(ctx context.Context, taskId string) error {
					return nil
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(eservice)
			controller.DeleteTask(c)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("error deleting task", func() {
			eservice := services.MockTaskService{
				FakeDeleteTaskById: func(ctx context.Context, taskId string) error {
					return fmt.Errorf("failed to delete task")
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			controller := NewTaskController(eservice)
			controller.DeleteTask(c)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Failed to delete task"))
		})

		It("missing task id", func() {
			eservice := services.MockTaskService{}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/tasks/", nil)
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: ""}} // Missing task id

			controller := NewTaskController(eservice)
			controller.DeleteTask(c)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var response *commons.ApiErrorResponsePayload
			uerr := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(uerr).NotTo(HaveOccurred())
			Expect(response.Status).To(Equal("Error"))
			Expect(response.Message).To(Equal("Task ID is required"))
		})
	})
})
