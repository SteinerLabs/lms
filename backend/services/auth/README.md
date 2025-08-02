# Auth Service

The Auth Service is responsible for user authentication, authorization, and identity management within the LMS platform. This service handles user login, registration, permission management, and secure access to other microservices.

## Features

- User management (create, read, update, delete)
- Authentication (login, logout, token validation)
- Role-based access control
- JWT token generation and validation
- Multi-factor authentication support
- Password management (hashing, reset)
- Event publishing for auth events

## Architecture

The Auth Service is built using a layered architecture:

- **API Layer**: Handles incoming requests and outgoing responses
- **Service Layer**: Contains the business logic
- **Repository Layer**: Handles data persistence
- **Event Layer**: Publishes events to Kafka

## Data Models

The Auth Service uses the following data models:

- **User**: Represents a user in the system
- **Role**: Represents a set of permissions
- **Permission**: Represents a specific action that can be performed
- **UserRole**: Represents the many-to-many relationship between users and roles
- **RolePermission**: Represents the many-to-many relationship between roles and permissions
- **Session**: Represents a user's authenticated session
- **MFADevice**: Represents a multi-factor authentication device
- **PasswordReset**: Represents a password reset request

## Database

The Auth Service uses PostgreSQL for data storage. The schema includes tables for users, roles, permissions, user_roles, role_permissions, sessions, mfa_devices, and password_resets.

## Events

The Auth Service publishes the following events:

- **user.created**: When a new user is created
- **user.updated**: When a user's information is updated
- **user.deleted**: When a user is deleted
- **user.login**: When a user logs in
- **user.logout**: When a user logs out

## API

The Auth Service provides a gRPC API for other services to use. The API includes methods for:

- User management
- Authentication
- Role management

## Configuration

The Auth Service can be configured using environment variables:

- **SERVER_PORT**: The port the server listens on (default: 50051)
- **DB_HOST**: The database host (default: localhost)
- **DB_PORT**: The database port (default: 5432)
- **DB_USER**: The database user (default: postgres)
- **DB_PASSWORD**: The database password (default: postgres)
- **DB_NAME**: The database name (default: auth)
- **DB_SSL_MODE**: The database SSL mode (default: disable)
- **JWT_SECRET**: The JWT secret key (default: your-secret-key)
- **JWT_ACCESS_TOKEN_TTL**: The JWT access token TTL in minutes (default: 15)
- **JWT_REFRESH_TOKEN_TTL**: The JWT refresh token TTL in days (default: 7)
- **JWT_ISSUER**: The JWT issuer (default: lms-auth-service)
- **JWT_AUDIENCE**: The JWT audience (default: lms-api)
- **KAFKA_BROKERS**: The Kafka brokers (default: localhost:9092)
- **KAFKA_TOPIC**: The Kafka topic (default: user-events)

## Getting Started

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 13 or higher
- Kafka (optional, for event publishing)

### Installation

1. Clone the repository
2. Navigate to the auth service directory
3. Build the service:

```bash
go build -o auth-service ./cmd
```

### Running the Service

```bash
./auth-service
```

### Running the Tests

```bash
go test ./...
```

## Usage

### Creating a User

```go
user, err := authService.CreateUser(ctx, "user@example.com", "password123", "John", "Doe")
if err != nil {
    // Handle error
}
```

### Authenticating a User

```go
session, err := authService.Login(ctx, "user@example.com", "password123", "127.0.0.1", "Mozilla/5.0")
if err != nil {
    // Handle error
}

// Use the session token for authenticated requests
token := session.Token
```

### Validating a Token

```go
userID, permissions, err := authService.ValidateToken(ctx, token)
if err != nil {
    // Handle error
}

// Check if the user has the required permissions
hasPermission := false
for _, permission := range permissions {
    if permission == "resource:action" {
        hasPermission = true
        break
    }
}
```

### Creating a Role

```go
role, err := authService.CreateRole(ctx, "admin", "Administrator with full access")
if err != nil {
    // Handle error
}
```

### Assigning a Role to a User

```go
err := authService.AssignRoleToUser(ctx, userID, roleID)
if err != nil {
    // Handle error
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.