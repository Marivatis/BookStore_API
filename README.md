# Book Store API
A backend service for managing books, magazines, and customer orders. Written in Go.

This project was built to practice manual SQL handling and structuring basic domain logic.
The code is split into layers (entities, DTOs, services, repositories) and uses manual SQL with transaction handling in key operations.
Configuration is managed via environment variables. All logging is structured with zap.

## Features
- Full CRUD for books, magazines, and orders
- Manual SQL queries using pgx
- Transactional operations
- Structured logging with zap
- Environment-based configuration
- PostgreSQL for persistent storage
- Database migrations using golang-migrate

## Technologies
- Go 1.24
- Echo web framework
- PostgreSQL
- Docker & Docker Compose
- golang-migrate
- zap logging

## API Endpoints

Below is a basic overview of the available endpoints for the API.  
Book endpoints are shown with an example request body for creating a new book.  
Other entity endpoints follow a similar pattern.
You can also explore and test the API using [this Postman collection](https://marivatis-2453845.postman.co/workspace/Martyniuk-Ivan's-Workspace~a2fea552-0d21-4ce8-8d28-e4a009bbb42c/collection/47003959-8f3bfb68-c3d6-41cd-85b4-e33c7211286c?action=share&creator=47003959).

### Books
| Method | Path       | Description             |
|--------|------------|-------------------------|
| GET    | /books/:id | Get book by ID          |
| POST   | /books     | Create a new book       |
| PUT    | /books/:id | Update an existing book |
| DELETE | /books/:id | Delete a book by ID     |

Example: Create Book Request Body
```json
{
  "name": "Refactoring",
  "price": 38.50,
  "stock": 10,
  "author": "Martin Fowler",
  "isbn": "978-0201485677"
}
```

## How to run
Run locally with Go:
```bash
APP_ENV=development go run cmd/main.go
```  
Or with Make:
```bash
make run
```
Run with Docker (uses .env.docker file parameter):
```bash
docker compose --env-file .env.docker up --build
```
