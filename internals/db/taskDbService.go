package db

import (
	"context"
	"fmt"

	"TaskSvc/commons/appdb"
	"TaskSvc/configs"
	models "TaskSvc/internals/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbService struct {
	collection appdb.DatabaseCollection
}

type DbService interface {
	GetTaskById(context context.Context, taskId string) (*models.TaskSchema, error)
	SaveTask(context context.Context, task *models.TaskSchema) (string, error)
	UpdateTask(context context.Context, task *models.TaskSchema, taskId string) error
	DeleteTaskById(context context.Context, taskId string) error
	GetTasks(context context.Context) ([]*models.TaskSchema, error)
}

func NewDbService(dbclient appdb.DatabaseClient) DbService {
	return &dbService{
		collection: dbclient.Collection(configs.MONGO_TASK_COLLECTION),
	}
}

func (d *dbService) GetTaskById(ctx context.Context, taskId string) (*models.TaskSchema, error) {
	var task models.TaskSchema
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, fmt.Errorf("invalid taskId: %v", err)
	}
	err = d.collection.FindOne(ctx, bson.M{"_id": id}, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (d *dbService) GetTasks(ctx context.Context) ([]*models.TaskSchema, error) {
	var tasks []*models.TaskSchema
	err := d.collection.Find(ctx, bson.M{}, &options.FindOptions{}, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %v", err)
	}
	return tasks, nil
}

func (d *dbService) SaveTask(ctx context.Context, task *models.TaskSchema) (string, error) {
	result, err := d.collection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}
	taskID := result.InsertedID.(primitive.ObjectID).Hex()
	return taskID, nil
}

func (d *dbService) UpdateTask(ctx context.Context, task *models.TaskSchema, taskId string) error {
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return fmt.Errorf("invalid taskId: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"updatedAt":   task.UpdatedAt,
		},
	}

	_, err = d.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbService) DeleteTaskById(ctx context.Context, taskId string) error {
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return fmt.Errorf("invalid taskId: %v", err)
	}
	_, err = d.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
