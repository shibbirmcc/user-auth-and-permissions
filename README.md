# **User Auth and Permissions Service**

[![Go Version](https://img.shields.io/badge/Go-1.22.8-blue)](https://golang.org)
[![Build Status](https://github.com/shibbirmcc/user-auth-and-permissions/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/shibbirmcc/user-auth-and-permissions/actions)
[![Coverage Status](https://codecov.io/gh/shibbirmcc/user-auth-and-permissions/branch/develop/graph/badge.svg)](https://codecov.io/gh/shibbirmcc/user-auth-and-permissions)



[![License](https://img.shields.io/github/license/shibbirmcc/user-auth-and-permissions)](LICENSE)

## **Overview**

The **User Auth and Permissions** service is a robust, scalable **Golang** backend microservice for managing **user authentication**, **role-based access control (RBAC)**, and **permissions management**. Built using **PostgreSQL**, this service ensures secure **user registration** with email confirmation, **JWT authentication**, and flexible **permission management** through **RESTful APIs**. It's designed for use in **microservices** architectures with detailed **API documentation** using **Swagger**.

## **Table of Contents**

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Database Setup](#database-setup)
  - [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Running Tests](#running-tests)
- [Code Coverage](#code-coverage)
- [Contributing](#contributing)
- [License](#license)

---

## **Features**

- **User Registration**: Includes email confirmation for secure user onboarding.
- **JWT Authentication**: Provides secure login, returning JWT tokens with `username`, `userId`, `role`, and `permissions`.
- **Role-Based Access Control (RBAC)**: Granular role and permission management.
- **Database Migrations**: Automatic **PostgreSQL** schema migrations handled at service startup.
- **Logging**: Structured logging with `logrus` for **info**, **warning**, and **error** levels.
- **Swagger API Documentation**: Automatically generated **OpenAPI** docs for easy integration.
- **Unit and Integration Testing**: Extensive **unit tests** and **integration tests** using **TestContainers**.

---

## **Tech Stack**

- **Language**: [Golang](https://golang.org)
- **Frameworks**:
  - [Gin](https://github.com/gin-gonic/gin) for routing and middleware.
  - [GORM](https://gorm.io/) for ORM and database interaction.
  - [Tonic](https://github.com/loopfz/golang-swiss-army-knife/tree/master/tonic) for parameter binding in **Gin**.
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Authentication**: [JWT](https://github.com/golang-jwt/jwt)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Testing**: [TestContainers](https://github.com/testcontainers/testcontainers-go) for integration testing.

---

## **Getting Started**

### **Prerequisites**

To run this **Golang microservice** locally, ensure you have the following installed:

- **Golang**: Version 1.18 or higher.
- **PostgreSQL**: Latest version.
- **Docker**: Required for running PostgreSQL locally during tests.

### **Installation**

1. Clone the repository:

   ```bash
   git clone https://github.com/shibbirmcc/user-auth-and-permissions.git
   cd user-auth-and-permissions
   ```
2. Install The Dependencies:
    ```bash
    go mod download
    ```
3. Database Setup:
    - Install or run postgresql database. If using docker then run these commands:
    ```bash
    docker pull postgres
    # replace the myuser with your desired admin username and mypassword with your desired admin password
    docker run --name my-postgres-container -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=mydatabase -p 5432:5432 -d postgres
    ```
    - Install the postgresql-client to execute db operations from terminal:
    ```bash
    sudo apt install postgresql-client

    ```
    - Connect to the default Database from terminal:
    ```bash
    psql -h localhost -p 5432 -U myuser
    ```
    - Create a seperate user credentials for the service to use for db operations:
    ```bash
    # replace serviceuser and servicepassword with your desired value
    CREATE ROLE serviceuser WITH LOGIN PASSWORD 'servicepassword';
    ALTER ROLE serviceuser CREATEDB;
    ```
    - Create a new database for the service:
    ```bash
    # replace servicedatabase with desired value
    CREATE DATABASE servicedatabase OWNER serviceuser;
    ALTER DATABASE servicedatabase OWNER TO serviceuser;
    ```
    - Grant Privileges to the New User:
    ```bash
    GRANT ALL PRIVILEGES ON DATABASE servicedatabase TO serviceuser;
    ```
    - Put the service user and password in the `.env` file:
    ```bash
    DB_HOST=localhost
    DB_USER=yourServiceUser
    DB_PASSWORD=yourServicePassword
    DB_NAME=yourServiceDatabase
    DB_PORT=5432
    ```
4. Running the Application:
    ```bash
    go run main.go
    ```

### **API Documentation**
This service comes with Swagger documentation for all APIs. After running the service, you can access the API docs at:
```bash
http://localhost:8080/swagger/index.html
```
The following endpoints are available:

* POST /register: Register a new user with email confirmation.
* POST /login: Authenticate the user and get JWT tokens.
* GET /roles: Fetch available roles.
* POST /roles: Add new roles (Admin only).
* GET /permissions: Fetch available permissions.

### **Running Tests**
This project includes both unit and integration tests using TestContainers to simulate the PostgreSQL database.

- **Clear the cache before running tests**:
```bash
go clean -testcache
```
- **Run Unit Tests**:
```bash
go test ./...
```
- **Run Integration Tests**:
```bash
go test -tags=integration ./tests
```
- **Show Test Execution time for each package**:
```bash
go test -v ./...
```

### **Code Coverage**
To check the code coverage for your tests, use the following commands:
- Run tests and generate coverage:
```bash
go test -coverprofile=coverage.out ./...
```
- View coverage
```bash
go tool cover -html=coverage.out
```

### **Contributing**
We welcome contributions to this project! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new feature branch.
3. Make your changes, adding tests where necessary.
4. Ensure that your code adheres to the existing coding conventions and passes all tests.
5. Submit a pull request.

## **License**
This project is licensed under the MIT License. See the [LICENSE](https://github.com/shibbirmcc/user-auth-and-permissions/blob/main/LICENSE) file for more details.

## **Contact**
For any inquiries or issues, please [open an issue](https://github.com/shibbirmcc/user-auth-and-permissions/issues) on GitHub.
