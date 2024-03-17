
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
| `ID`      | `int` | **Required**. Task ID|
| `Description`      | `string` | **Required**. Task Description|
| `Completed`      | `boolean` | **Required**. used to label task Completed |

#### Delete an existing task

```http
  DELETE /delete_task
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `ID`      | `int` | **Required**. Task ID |

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

specify with port number default is 8080

```bash
  go run main.go task.go helper.go -port=8081
```

Follow API reference for the current resources and Methods allowed

Test using CURL

``` all_tasks
curl -X GET 'http://localhost:8080/all_tasks'

Response:

[{"id":1,"description":"create a Go web API app for my portfolio","completed":false},{"id":2,"description":"Get coffee with my parents","completed":false}]
```

``` add_task
curl -X POST \
  http://localhost:8080/add_task \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "new task",
    "completed": false
}'

Response:
{"status":"ok","message":"Task added successfully"}
```

``` update_task
curl -X POST \
  http://localhost:8080/update_task \
  -H 'Content-Type: application/json' \
  -d '{
    "ID": 1,
    "description": "create a Go web API app for my portfolio",
    "completed": true
}'
Response:
{"status":"ok","message":"Updated Task successfully"}

```

``` delete_tasl
curl -X DELETE \
  http://localhost:8080/delete_task \
  -H 'Content-Type: application/json' \
  -d '{
    "ID": 2 
}'
Response:
{"status":"ok","message":"Task deleted successfully"}
```