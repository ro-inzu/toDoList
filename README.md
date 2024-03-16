
# ToDo List

Keep track of all your todo items


## API Reference

#### Get all tasks

```http
  GET /get_all_tasks
```

#### Update task 

```http
  POST /update_task
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Description`      | `string` | **Required**. Task Description|
| `Completed`      | `boolean` | **Required**. used to label task Completed |

#### Delete an existing task

```http
  DELETE /delete_task
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of task to delete |

####  Add a new Task

```http
  POST /add_task
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Description`      | `string` | **Required**. Task Description|
| `Completed`      | `boolean` | **Required**. used to label task Completed |


## Features

- create new task
- delete existing task
- update existing task
- view all task


## Run Locally

Clone the project

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd toDoList
```

Install Go

```bash
  go run main.go helper.go task.go
```

Follow API reference for the current resources and Methods allowed
