# ZPE Cloud - User Management Service

## Description

This repository contains the implementation of a microservice for the platform "ZPE Cloud", designed as part of a technical challenge for a selection process. The microservice provides a comprehensive API for managing user access and roles within the system, catering to different user roles including Admin, Modifier, and Watcher. This service is responsible for managing users in the ZPE Cloud. It is a RESTful service that provides endpoints for creating, updating, deleting, and retrieving users.

## Structure

The microservice is built using Golang and is structured into several packages:

- `cmd/server`: Contains the main entry point for the application.
- `config`: Handles configuration loading.
- `internal/user`: Contains the business logic for user management, including handlers, models, and storage.
- `internal/msgs`: Contains response messages.
- `scripts`: Contains scripts for setting project execution.

### Endpoints

#### Create User
- **POST** `/users`
  - **Headers:** `X-User-Type: <role>`
  - **Payload:** `{"name": "Han Solo", "email": "solo@example.com", "roles": ["Admin"]}`
  - **Response:**
    - `201 Created`: `{"id":"<user_id>", "message":"User created successfully"}`
    - `409 Conflict`: `{"message":"user already exists"}`
    - `403 Forbidden`: `{"message":"forbidden"}`

#### List Users
- **GET** `/users`
  - **Headers:** `X-User-Type: <role>`
  - **Response:**
    - `200 OK`: `[{"id":"<user_id>", "name":"<name>", "email":"<email>", "roles":["<role>"]}]`
    - `403 Forbidden`: `{"message":"forbidden"}`

#### Get User Details
- **GET** `/users/{id}`
  - **Headers:** `X-User-Type: <role>`
  - **Response:**
    - `200 OK`: `{"id":"<user_id>", "name":"<name>", "email":"<email>", "roles":["<role>"]}`
    - `404 Not Found`: `[]`
    - `403 Forbidden`: `{"message":"forbidden"}`

#### Update User Roles
- **PUT** `/users/roles/{id}`
  - **Headers:** `X-User-Type: <role>`
  - **Payload:** `{"roles": ["Modifier"]}`
  - **Response:**
    - `200 OK`: `{"message":"User roles updated successfully"}`
    - `403 Forbidden`: `{"message":"forbidden"}`
    - `404 Not Found`: `{"message":"user not found"}`

#### Delete User
- **DELETE** `/users/{id}`
  - **Headers:** `X-User-Type: <role>`
  - **Response:**
    - `204 No Content`
    - `404 Not Found`: `{"message":"user not found"}`
    - `403 Forbidden`: `{"message":"forbidden"}`

### Usage Examples

### Usage Examples

#### Create a User
```sh
curl -i -X POST -H "X-User-Type: Admin" -H "Content-Type: application/json" -d '{"name":"Han Solo","roles":["Admin"],"email":"solo@example.com"}' http://localhost:8080/users
```

#### List Users
```sh
curl -i -X GET -H "X-User-Type: Admin" http://localhost:8080/users
```

#### Get User Details
```sh
curl -i -X GET -H "X-User-Type: Admin" http://localhost:8080/users/1
```

#### Update User Roles
```sh
curl -i -X PUT -H "X-User-Type: Admin" -H "Content-Type: application/json" -d '{"roles":["Modifier"]}' http://localhost:8080/users/roles/1
```

#### Delete a User
```sh
curl -i -X DELETE -H "X-User-Type: Admin" http://localhost:8080/users/1
```

#### Handle Unauthorized Actions
```sh
curl -i -X POST -H "X-User-Type: Watcher" -H "Content-Type: application/json" -d '{"name":"Darth Vader","roles":["Admin"],"email":"vader@example.com"}' http://localhost:8080/users
```

### Setup

#### Prerequisites
- Go 1.22 or later
- Docker

#### Running Locally
1. Clone the repository:
    ```bash
    git clone https://github.com/dahn94/zpe-cloud-user-management-service.git
    cd zpe-cloud-user-management-service
    ```

2. Run the application locally using the following command:
    ```bash
    go run cmd/server/main.go
    ```

#### Running with Docker
1. Build and run the Docker image:
    ```bash
    docker build -t zpe-cloud-user-management-service .
    docker run --env-file .env -p 8080:8080 zpe-cloud-user-management-service

    ```

### Testing
Run the tests using the following command:
```bash
go test -v ./internal/user
```