# Auth Service Data Models

## Overview

The Auth Service is responsible for user authentication, authorization, and identity management. This document defines the core data models used by the Auth Service.

## Data Models

### User

The User model represents the authentication and authorization information for a user.

```go
type User struct {
    ID              string    `json:"id" db:"id"`                           // Unique identifier
    Email           string    `json:"email" db:"email"`                     // User's email address
    PasswordHash    string    `json:"-" db:"password_hash"`                 // Hashed password
    FirstName       string    `json:"first_name" db:"first_name"`           // User's first name
    LastName        string    `json:"last_name" db:"last_name"`             // User's last name
    Active          bool      `json:"active" db:"active"`                   // Whether the user is active
    EmailVerified   bool      `json:"email_verified" db:"email_verified"`   // Whether the email is verified
    MFAEnabled      bool      `json:"mfa_enabled" db:"mfa_enabled"`         // Whether MFA is enabled
    LastLogin       time.Time `json:"last_login" db:"last_login"`           // Last login timestamp
    FailedAttempts  int       `json:"failed_attempts" db:"failed_attempts"` // Number of failed login attempts
    Locked          bool      `json:"locked" db:"locked"`                   // Whether the account is locked
    LockExpiry      time.Time `json:"lock_expiry" db:"lock_expiry"`         // When the lock expires
    CreatedAt       time.Time `json:"created_at" db:"created_at"`           // Creation timestamp
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`           // Last update timestamp
}
```

### Role

The Role model represents a set of permissions that can be assigned to users.

```go
type Role struct {
    ID          string    `json:"id" db:"id"`               // Unique identifier
    Name        string    `json:"name" db:"name"`           // Role name
    Description string    `json:"description" db:"description"` // Role description
    CreatedAt   time.Time `json:"created_at" db:"created_at"`   // Creation timestamp
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // Last update timestamp
}
```

### Permission

The Permission model represents a specific action that can be performed.

```go
type Permission struct {
    ID          string    `json:"id" db:"id"`               // Unique identifier
    Name        string    `json:"name" db:"name"`           // Permission name
    Description string    `json:"description" db:"description"` // Permission description
    Resource    string    `json:"resource" db:"resource"`       // Resource the permission applies to
    Action      string    `json:"action" db:"action"`           // Action (create, read, update, delete)
    CreatedAt   time.Time `json:"created_at" db:"created_at"`   // Creation timestamp
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // Last update timestamp
}
```

### UserRole

The UserRole model represents the many-to-many relationship between users and roles.

```go
type UserRole struct {
    ID        string    `json:"id" db:"id"`             // Unique identifier
    UserID    string    `json:"user_id" db:"user_id"`       // User ID
    RoleID    string    `json:"role_id" db:"role_id"`       // Role ID
    CreatedAt time.Time `json:"created_at" db:"created_at"` // Creation timestamp
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update timestamp
}
```

### RolePermission

The RolePermission model represents the many-to-many relationship between roles and permissions.

```go
type RolePermission struct {
    ID           string    `json:"id" db:"id"`                   // Unique identifier
    RoleID       string    `json:"role_id" db:"role_id"`           // Role ID
    PermissionID string    `json:"permission_id" db:"permission_id"` // Permission ID
    CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Creation timestamp
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`       // Last update timestamp
}
```

### Session

The Session model represents a user's authenticated session.

```go
type Session struct {
    ID           string    `json:"id" db:"id"`                   // Unique identifier
    UserID       string    `json:"user_id" db:"user_id"`           // User ID
    Token        string    `json:"token" db:"token"`               // Session token
    RefreshToken string    `json:"refresh_token" db:"refresh_token"` // Refresh token
    ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`       // Expiration timestamp
    IP           string    `json:"ip" db:"ip"`                     // IP address
    UserAgent    string    `json:"user_agent" db:"user_agent"`       // User agent
    CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Creation timestamp
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`       // Last update timestamp
}
```

### MFADevice

The MFADevice model represents a multi-factor authentication device.

```go
type MFADevice struct {
    ID        string    `json:"id" db:"id"`             // Unique identifier
    UserID    string    `json:"user_id" db:"user_id"`       // User ID
    Type      string    `json:"type" db:"type"`           // Device type (app, sms, email)
    Secret    string    `json:"-" db:"secret"`           // Secret key
    Verified  bool      `json:"verified" db:"verified"`     // Whether the device is verified
    CreatedAt time.Time `json:"created_at" db:"created_at"` // Creation timestamp
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update timestamp
}
```

### PasswordReset

The PasswordReset model represents a password reset request.

```go
type PasswordReset struct {
    ID        string    `json:"id" db:"id"`             // Unique identifier
    UserID    string    `json:"user_id" db:"user_id"`       // User ID
    Token     string    `json:"token" db:"token"`         // Reset token
    ExpiresAt time.Time `json:"expires_at" db:"expires_at"` // Expiration timestamp
    Used      bool      `json:"used" db:"used"`           // Whether the token has been used
    CreatedAt time.Time `json:"created_at" db:"created_at"` // Creation timestamp
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update timestamp
}
```

## Relationships

- A User can have multiple Roles through UserRole
- A Role can have multiple Permissions through RolePermission
- A User can have multiple Sessions
- A User can have multiple MFADevices
- A User can have multiple PasswordResets

## Database Schema

The Auth Service uses PostgreSQL for data storage. The schema includes the following tables:

- users
- roles
- permissions
- user_roles
- role_permissions
- sessions
- mfa_devices
- password_resets

## Events

The Auth Service publishes the following events:

### UserCreated

```go
type UserCreated struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### UserUpdated

```go
type UserUpdated struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### UserDeleted

```go
type UserDeleted struct {
    ID        string    `json:"id"`
    DeletedAt time.Time `json:"deleted_at"`
}
```

### UserLoggedIn

```go
type UserLoggedIn struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    IP        string    `json:"ip"`
    UserAgent string    `json:"user_agent"`
    LoginAt   time.Time `json:"login_at"`
}
```

### UserLoggedOut

```go
type UserLoggedOut struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    LogoutAt  time.Time `json:"logout_at"`
}
```