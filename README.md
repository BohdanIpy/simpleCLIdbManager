## Simple CLI User Management System

A **lightweight Command-Line Interface (CLI)** application written in **Go** that manages users stored in a **PostgreSQL** database.  
It supports user creation, retrieval, updates, and deletions — all directly from the terminal.

---



This project demonstrates a clean, modular design for building database-driven Go applications.  
It includes:
- Repository pattern for clean data access
- Decorator pattern for internal modules and middleware 
- Middleware for logging and transactional handling
- Unit tests using [sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- Docker setup for PostgreSQL

---

## Features

- CRUD operations for users  
- Transactional middleware for safe DB operations  
- Logging middleware to track queries  
- Fully testable SQL layer using mocks  
- Docker Compose integration for easy database setup  

---

## Tech Stack

| Component | Description |
|------------|-------------|
| **Go** | Core application language |
| **PostgreSQL** | Database |
| **sqlmock** | Mocking framework for SQL tests |
| **Docker Compose** | Local database setup |

---

# Example CLI Flow
----------------------------

Choose an action:
1. Get all users
2. Get user by ID
3. Add new user
4. Update user
5. Delete user
6. Exit

> 1
<br>
{"John", "john@example.com", "123456", "datetime"}


## Installation & Setup

### Clone the repository
```bash
git clone https://github.com/YourUsername/simpleCLIdbManager.git
cd simpleCLIdbManager

(before run docker compose up -d)
go run cmd/cliManager/main.go host post password database 
```

Tests located in the internal/middleware and internal/repository/impl
