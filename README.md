# Task Management Golang Application

## Overview

This is a Task Management application developed in Go. It provides a set of APIs for managing tasks. The application runs locally and connects to a MongoDB database.

## Run Locally

### Clone the project

```bash
git clone https://github.com/pradeep-thombre/Task-Management.git
```

### Go to the project directory

```bash
cd Task-Management
```

### Import dependencies

```bash
go mod tidy
```

### Start the server

```bash
go run .
```

## API Reference

### Get all Tasks

```http
GET /tasks
```

Gets a list of all tasks and the total count of tasks present.

### Get Task by Id

```http
GET /tasks/${id}
```

| Parameter | Type     | Description                        |
| :-------- | :------- | :--------------------------------- |
| `id`      | `string` | **Required**. ID of task to fetch  |

Gets the task by the provided ID.

### Delete Task by Id

```http
DELETE /tasks/${id}
```

| Parameter | Type     | Description                        |
| :-------- | :------- | :--------------------------------- |
| `id`      | `string` | **Required**. ID of task to delete |

Deletes the task by the provided ID.

### Create a new Task

```http
POST /tasks
```

Payload:
```json
{
    "title": "string",        // required
    "description": "string",  // required
    "status": "string"        // required
}
```

Creates a new task with the provided payload and returns the ID of the task.

### Update Task by Id

```http
PATCH /tasks/${id}
```

| Parameter | Type     | Description                        |
| :-------- | :------- | :--------------------------------- |
| `id`      | `string` | **Required**. ID of task to update |

Payload:
```json
{
    "title": "string",        // required
    "description": "string",  // required
    "status": "string"        // required
}
```

Updates the task details by the provided ID and payload.


## Token

```
Bearer <your Token>
```


## Testing

```
ginkgo -r -v
```

## Authors

- [Pradeep Thombre](https://www.github.com/Pradeep-Thombre)

## 🛠 Tech Stacks

- Golang
- MongoDB
- Gin Framework
- Ginkgo
- Gomega


## Support

For support, email us at [pradeepbthombre@gmail.com](mailto:pradeepbthombre@gmail.com)
```
