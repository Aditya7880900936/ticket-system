# Ticket System

A small backend service for managing support tickets, built with Go, Gin, PostgreSQL, GORM, and JWT authentication.

## Tech Stack

- Go
- Gin
- PostgreSQL
- GORM
- JWT
- bcrypt
- Docker
- Docker Compose

## Features

- User registration
- User login with JWT authentication
- Password hashing using bcrypt
- Create tickets
- View own tickets
- View a specific own ticket
- Update ticket status
- Ownership-based authorization
- Ticket status transition validation
- PostgreSQL persistence
- Dockerized application
- Docker Compose setup with PostgreSQL

## Project Structure

```text
ticket-system/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   └── service/
├── .env.example
├── .gitignore
├── .dockerignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## API Endpoints

### Health

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### Register

```http
POST /auth/register
```

Request:

```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

### Login

```http
POST /auth/login
```

Request:

```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "<jwt-token>"
}
```

### Create Ticket

```http
POST /tickets
Authorization: Bearer <token>
```

Request:

```json
{
  "title": "Unable to login",
  "description": "Login is failing with an authentication error."
}
```

### List Own Tickets

```http
GET /tickets
Authorization: Bearer <token>
```

### Get Own Ticket

```http
GET /tickets/{id}
Authorization: Bearer <token>
```

### Update Ticket Status

```http
PATCH /tickets/{id}/status
Authorization: Bearer <token>
```

Request:

```json
{
  "status": "in_progress"
}
```

Allowed statuses:

```text
open
in_progress
closed
```

## Ticket Status Rules

Tickets follow a simple status flow:

```text
open → in_progress → closed
```

Rules:

- A newly created ticket starts as `open`.
- An `open` ticket can move to `in_progress`.
- An `open` ticket cannot directly move to `closed`.
- An `in_progress` ticket can move to `closed`.
- An `in_progress` ticket cannot move back to `open`.
- A `closed` ticket cannot be reopened or modified.

## Authorization

All ticket endpoints require a valid JWT.

Users can only access tickets they created.

Ownership is enforced when retrieving or updating a ticket by checking both:

```text
ticket_id
user_id
```

A user attempting to access another user's ticket receives:

```json
{
  "error": "ticket not found"
}
```

## Environment Variables

Create a `.env` file using `.env.example` as a reference.

Example:

```env
PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ticket_system

JWT_SECRET=your_jwt_secret
```

For production deployment, use a strong randomly generated JWT secret.

## Running Locally

### Without Docker

Make sure PostgreSQL is running and the environment variables are configured.

Then:

```bash
go mod download
go run ./cmd/server
```

The server will start on:

```text
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/health
```

## Running with Docker Compose

Build and start the application:

```bash
docker compose up --build
```

Or run in detached mode:

```bash
docker compose up -d --build
```

Check running containers:

```bash
docker compose ps
```

Health check:

```bash
curl http://localhost:8080/health
```

Expected:

```json
{
  "status": "ok"
}
```

Stop the application:

```bash
docker compose down
```

## Running with Docker

Build the image:

```bash
docker build -t ticket-system .
```

Run the container:

```bash
docker run -p 8080:8080 ticket-system
```

The application will be available at:

```text
http://localhost:8080
```

Note: when running the application container directly, PostgreSQL must be separately available and the appropriate database environment variables must be provided.

## Deployment

Deployment platform:

```text
Render
```

Deployed URL:

```text
TODO: Add deployed URL after deployment
```

Health endpoint:

```text
TODO: Add deployed health URL after deployment
```

Example:

```text
https://your-app.onrender.com/health
```

## Database

The application uses PostgreSQL with GORM.

The database schema contains two main tables:

### Users

- `id`
- `email`
- `password_hash`
- `created_at`
- `updated_at`

### Tickets

- `id`
- `user_id`
- `title`
- `description`
- `status`
- `created_at`
- `updated_at`

GORM AutoMigrate is used to create or update the required tables when the application starts.

## Assumptions

- There is no admin role in the system.
- Every ticket belongs to exactly one user.
- Users can only view and update their own tickets.
- Ticket status starts as `open`.
- Closed tickets cannot be reopened.
- No ticket assignment workflow is implemented.
- No comments or attachments are supported.
- JWT tokens expire after 24 hours.
- PostgreSQL is used for persistent storage.
- The implementation intentionally keeps the architecture simple and focused on the assignment requirements.

## HTTP Status Codes

| Status | Meaning |
| --- | --- |
| 200 | Successful request |
| 201 | Resource created |
| 400 | Invalid request |
| 401 | Missing or invalid authentication |
| 404 | Resource not found |
| 409 | Resource already exists |
| 500 | Internal server error |

## Architecture

The application follows a simple layered structure:

```text
HTTP Request
     ↓
Gin Router
     ↓
JWT Middleware
     ↓
Handler
     ↓
Repository
     ↓
PostgreSQL
```

Authentication-related password hashing and JWT operations are handled through the service layer.

## Testing

The API was manually tested for:

- User registration
- User login
- JWT authentication
- Ticket creation
- Ticket listing
- Ticket retrieval
- Ticket status updates
- Invalid status transitions
- Closed ticket protection
- Missing authentication
- Invalid ticket IDs
- Cross-user ticket access
- Cross-user ticket modification
- Dockerized application startup
- PostgreSQL connectivity
- Database migration
- Health endpoint

## License

This project was created as part of a backend development internship assignment.