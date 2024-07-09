# Integra Coding Assessment

This repository contains the source code for the Integra Coding Assessment project, which includes a Go backend and an Angular frontend. The project uses Docker to manage dependencies and run services in a containerized environment.

## Project Structure

```sh
integra-coding-assessment/
├── .gitignore
├── .env
├── go-backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── db/
│   │   ├── db.go
│   │   ├── migrations/
│   │   │   ├── 001_create_users_table.up.sql
│   │   │   ├── 001_create_users_table.down.sql
│   ├── handlers/
│   │   └── user_handler.go
│   ├── models/
│   │   └── user.go
│   └── docs/
│       └── swagger.go
├── angular-frontend/
│   ├── Dockerfile
│   ├── src/
│   ├── package.json
│   └── ...
└── docker-compose.yml
```

## Prerequisites

Ensure you have the following installed on your system:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Setting Up the Project

1. **Clone the Repository**:

    ```sh
    git clone https://github.com/sedmo/integra-coding-assessment.git
    cd integra-coding-assessment
    ```

2. **Create the `.env` File**:

    - In the root directory of your project, create a `.env` file:

      ```env
      POSTGRES_USER=yourusername
      POSTGRES_PASSWORD=yourpassword
      POSTGRES_DB=yourdatabase
      DATABASE_URL=postgres://yourusername:yourpassword@postgres:5432/yourdatabase?sslmode=disable
      API_URL=http://localhost:1323
      ```

3. **Build and Run the Docker Containers**:

    ```sh
    docker-compose up --build -d
    ```

4. **Verify the Setup**:

    Ensure all containers are running:

    ```sh
    docker-compose ps
    ```

## Services

### Go Backend

The Go backend runs on port `1323` and connects to a PostgreSQL database.

- **Dockerfile**: `go-backend/Dockerfile`
- **Service Definition**: `docker-compose.yml` under `backend` service

### Angular Frontend

The Angular frontend is built using Node.js and served using Nginx. It runs on port `80`.

- **Dockerfile**: `angular-frontend/Dockerfile`
- **Service Definition**: `docker-compose.yml` under `frontend` service

### PostgreSQL Database

A PostgreSQL database is used for data storage. It runs on port `5432`.

- **Service Definition**: `docker-compose.yml` under `postgres` service
- **Volume**: `postgres_data` for data persistence

API Documentation with Swagger
Swagger is used to generate and serve API documentation for the Go backend.

Swagger UI: <http://localhost:1323/swagger/index.html>

To generate Swagger documentation, use the swag CLI tool:

Install swag CLI:

```sh
go get -u github.com/swaggo/swag/cmd/swag
```

Generate Swagger Documentation:

```sh
swag init -g main.go
```

Running Migrations
Applying Migrations
Migrations are used to manage database schema changes.

Up Migrations:

Migrations are applied automatically when the backend service starts.

Down Migrations:

To roll back the latest migration:

```sh
docker-compose run backend ./main -migrate-down
```

Running Tests
Go Backend Tests
To run the tests for the Go backend:

Navigate to the go-backend Directory:

```sh
cd go-backend
```

Run the Tests:

```sh
go test ./...
```

Angular Frontend Tests
To run the tests for the Angular frontend:

Navigate to the angular-frontend Directory:

```sh
cd angular-frontend
```

Install Dependencies:

```sh
npm install
```

Run the Tests:

```sh
npm test
```

Accessing the Application
Backend API: <http://localhost:1323>
Frontend Application: <http://localhost>

Stopping the Application
To stop the running Docker containers:

```sh
docker-compose down
```

Cleaning Up
To remove all Docker containers, networks, and volumes associated with the project:

```sh
docker-compose down -v
```
