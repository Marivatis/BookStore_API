# Book Store API
A backend service for managing books and customer orders (CRUD API), written in Go.
This project demonstrates clean architecture principles, including clear separation
of concerns between entities, services, repositories, and DTOs, as well as
integration with PostgreSQL, migrations, and Docker.

## Status
This project is currently under development.

## Features
- Create, read, update, and delete books
- Create and manage orders
- PostgreSQL as persistent storage
- Database migrations with golang-migrate
- Environment-based configuration
- Unit and integration tests

## Technologies
- Go 1.24
- Echo web framework
- PostgreSQL
- Docker & Docker Compose
- golang-migrate
- Testify

## API Endpoints
### Books
TBC
### Orders
TBC

## How to run
Run locally:
```bash
go run cmd/main.go
```  
With Make:
```bash
make run
```
With Docker:
TBC

[//]: # (```bash)
[//]: # (docker-compose up --build)
[//]: # (```)