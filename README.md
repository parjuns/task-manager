# Task Manager API

This API allows users to manage tasks, register, log in, and perform CRUD operations on tasks.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Parjun2000/task-manager.git
   cd task-manager
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Set up environment variables:**

   - Create a `app.env` file with your environment variables.
   - Add or change below mention environment varibles in the file:  
      `DB_HOST`=localhost  
      `DB_PORT`=5432  
      `DB_USER`=root  
      `DB_PASSWORD`=secret  
      `DB_NAME`=task_manager_db  
      `SERVER_PORT`=8080  
      `JWT_KEY`=my-secret-key
   - Logs are generated in `app.log` file.

4. **Set up Database PostgreSQL in docker:**

   - Run command mention in makefile.

     ```bash
         make postgres

         make createdb

         make migrateup
     ```

   - This commands run the database in docker with the tables created in them.
   - Database schema is executed using go migrate cli.
   - `Migrations` dir contains up and down files for database.

## Running the API

Ensure database is up and then Run the API using:

```bash
go run main.go
```

## Endpoints

### Authentication

- `POST /api/v1/auth/register`: Register a new user.
- `POST /api/v1/auth/login`: Log in an existing user.

### Task Management

- `GET /api/v1/tasks`: Get all tasks.
- `POST /api/v1/tasks`: Create a new task.
- `GET /api/v1/tasks/{id}`: Get a task by ID.
- `PUT /api/v1/tasks/{id}`: Update a task by ID.
- `DELETE /api/v1/tasks/{id}`: Delete a task by ID.
- `POST /api/v1/tasks/mark-done`: Mark tasks as 'done' concurrently.

### Pagination, Sorting & Filtering

Use tasks endpoint with query params in the API, for pagination, sorting, & filtering.

- `GET /api/v1/tasks/?page=1&limit=5&status=done&sort_by=created_at&order=asc`:

  `page`: Allows to paginate through the task list.  
   `limit`: Allows to set limit per page for task list.  
   `status`: Allows to filter based on status of task in list (e.g., "todo," "in progress," "done").  
   `sort_by`: Allows to sort based on title, status, description.  
   `order`: Allows to order with ASC or DESC.

## API Documentation

Access the API documentation using Swagger:

- `Swagger UI`: [http://localhost:8080/api/v1swagger/index.html](http://localhost:8080/api/v1/swagger/index.html)
- `API Swagger JSON`: [http://localhost:8080/api/v1/swagger/doc.json](http://localhost:8080/api/v1/swagger/doc.json)

## Authentication & Authorization

- Implemented user registration and login functionality.
- Users have to get JWT token from login endpoint and use for task enpoints.
- Users are able to create, read, update, and delete tasks only if authenticated.

## Middleware for Authentication, Database, Logging & ErrorHandling

- Implemented middleware for authentication of every endpoints.
- Implemented middleware for database with Go interface so that mock database can also be integrated for testing.
- Implemented middleware for logging each incoming request on console as well as in a log file `app.log`.
- ErrorHandling is done for all endpoints in middleware as well as handlers and provided meaningful error responses with http status codes.
- Implemented proper validation for task data (e.g., required fields, valid status values) is taken care of with validation packages.

## Concurrency

- `POST /api/v1/tasks/mark-done`: Mark tasks as 'done' concurrently.
- `Mark-done` route can be used which allows users to mark multiple tasks as "done" concurrently using Goroutines and channels for concurrent processing.

## Database

- `PostgreSQL's` relational data model is well-suited for structured data like tasks, ensuring data consistency and integrity, which is crucial for a task management system.

- The ability of `PostgreSQL` to handle complex queries, and relationships among data makes it suitable for task management, where task assignments, status updates, and dependencies might require intricate querying support.

  ### Schema for Users & Tasks Table

  | Column Name | Data Type    | Description             |
  | ----------- | ------------ | ----------------------- |
  | id          | INT          | Unique ID               |
  | username    | VARCHAR(50)  | Username                |
  | password    | VARCHAR(100) | Encrypted user password |

  | Column Name | Data Type    | Description                                                  |
  | ----------- | ------------ | ------------------------------------------------------------ |
  | task_id     | INT          | Unique ID                                                    |
  | title       | VARCHAR(100) | Title of the task                                            |
  | description | TEXT         | Description of the task                                      |
  | status      | VARCHAR(20)  | Task status (e.g., 'todo', 'in progress', 'done')            |
  | user_id     | INT          | Unique ID of the task owner (user_id referencing User table) |
  | created_at  | TIMESTAMP    | Date and time of creation                                    |

## Docker Containerize & Deploy on cloud platform

- `Dockerfile` & `docker-compose.yml` files contains docker image details and all necessary script run commands.
- Run the server with database using docker container with this command:
  ```bash
  docker compose up
  ```

## Deploy on AWS cloud.

- Replace localhost to the public IP of AWS-EC2 for endpoints:
- Access the API documentation deployed on aws using Swagger on Web browser:

  - `Swagger UI on AWS EC2`: [http://13.232.229.210:8080/api/v1swagger/index.html](http://13.232.229.210:8080/api/v1/swagger/index.html)

- Note: Public IP is used above just for illustration and should not use on production servers.

## Testing

- Handlers are tested using mock database.

- Perform unit tests on the API using:

  ```bash
  go test -v ./...
  ```
