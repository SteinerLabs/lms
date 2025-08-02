# Auth Service Data Models

## Overview

The Auth Service is responsible for user authentication, authorization, and identity management. This document defines the core data models used by the Auth Service.

## Data Models

### User

The User model represents the authentication and authorization information for a user.

```go
type AuthUser struct {
    ID              string // uuidv6 or ulid
    Email           string
    PasswordHash    string
    Active          bool
    EmailVerified   bool
    LastLogin       time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### RefreshToken

The RefreshToken model represents a refresh token for a user

```go
type RefreshToken struct {
	ID string
	UserId string
	TokenHash string // SHA-256 hash of the token
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt time.Time
	UserAgent string
	IpAddress string
}

// GenerateSecureToken returns a base64 URL-safe random string
func GenerateSecureToken(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// HashToken hashes a token using SHA256
func HashToken(token string) string {
    h := sha256.Sum256([]byte(token))
    return base64.RawURLEncoding.EncodeToString(h[:])
}

func CreateRefreshToken(db *sql.DB, userID string, userAgent, ip string) (string, error) {
    rawToken, err := GenerateSecureToken(32) // 256-bit
    if err != nil {
        return "", err
    }

    hashed := HashToken(rawToken)
    expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
    id := uuid.New().String()

    _, err = db.Exec(`
            INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, user_agent, ip_address)
            VALUES ($1, $2, $3, $4, NOW(), $5, $6)
        `, id, userID, hashed, expiresAt, userAgent, ip)
    if err != nil {
        return "", err
    }

    // Return plain token (not hashed) to client
    return rawToken, nil
}

func RotateRefreshToken(db *sql.DB, userID, providedToken string, userAgent, ip string) (string, string, error) {
    hashed := HashToken(providedToken)
    
    // Lookup existing token
    var tokenID string
    var expiresAt time.Time
    err := db.QueryRow(`
            SELECT id, expires_at FROM refresh_tokens
            WHERE user_id = $1 AND token_hash = $2 AND revoked_at IS NULL
        `, userID, hashed).Scan(&tokenID, &expiresAt)
    if err != nil {
        return "", "", fmt.Errorf("invalid refresh token")
    }
    
    if time.Now().After(expiresAt) {
        return "", "", fmt.Errorf("refresh token expired")
    }
    
    // Revoke old token
    _, err = db.Exec(`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`, tokenID)
    if err != nil {
        return "", "", err
    }
    
    // Create a new refresh token
    newRefreshToken, err := CreateRefreshToken(db, userID, userAgent, ip)
    if err != nil {
        return "", "", err
    }
    
    // Issue a new access token (JWT)
    newAccessToken, err := IssueAccessToken(userID) // implement JWT creation
    if err != nil {
        return "", "", err
    }
    
    return newAccessToken, newRefreshToken, nil
}

```

#### Flow

1. User logs in -> AuthService:
   - Issues short-lived JWT access token (e.g., 15min)
   - Generates a long-lived refresh token (random 256-bit value)
   - Hases and stores it in DB
2. When access token expires:
    - Client sends refresh token to AuthService
    - AuthService validates token (hash lookup and expiry/revokation)
    - Issues new access token (possibly a new refresh token with rotation)
    - Optionally invalidates old refresh token
3. On logout:
   - Delete or mark refresh token revoked in DB

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

- A User can have multiple PasswordResets

## Database Schema

The Auth Service uses PostgreSQL for data storage. The schema includes the following tables:

- password_resets

## Events

The Auth Service publishes the following events:

### UserCreated

```go
type AuthUserCreated struct {
    ID        string
    Email     string
	Username  string
	DisplayName string
    FirstName string
    LastName  string
	Bio string
	AvatarURL string
	Location string
	Links []Link
	Prefrences UserPrefrences
    CreatedAt time.Time
}
```

### AuthUserUpdated

```go
type AuthUserUpdated struct {
    ID string
    Email string
	Username string
    UpdatedAt time.Time
}
```

### AuthUserDeleted

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