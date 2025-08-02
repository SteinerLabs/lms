package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents the authentication and authorization information for a user
type User struct {
	ID              string    `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Active          bool      `json:"active" db:"active"`
	EmailVerified   bool      `json:"email_verified" db:"email_verified"`
	MFAEnabled      bool      `json:"mfa_enabled" db:"mfa_enabled"`
	LastLogin       time.Time `json:"last_login" db:"last_login"`
	FailedAttempts  int       `json:"failed_attempts" db:"failed_attempts"`
	Locked          bool      `json:"locked" db:"locked"`
	LockExpiry      time.Time `json:"lock_expiry" db:"lock_expiry"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user with default values
func NewUser(email, passwordHash, firstName, lastName string) *User {
	now := time.Now().UTC()
	return &User{
		ID:             uuid.New().String(),
		Email:          email,
		PasswordHash:   passwordHash,
		FirstName:      firstName,
		LastName:       lastName,
		Active:         true,
		EmailVerified:  false,
		MFAEnabled:     false,
		LastLogin:      time.Time{},
		FailedAttempts: 0,
		Locked:         false,
		LockExpiry:     time.Time{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Role represents a set of permissions that can be assigned to users
type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewRole creates a new role with default values
func NewRole(name, description string) *Role {
	now := time.Now().UTC()
	return &Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Permission represents a specific action that can be performed
type Permission struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewPermission creates a new permission with default values
func NewPermission(name, description, resource, action string) *Permission {
	now := time.Now().UTC()
	return &Permission{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	RoleID    string    `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewUserRole creates a new user role with default values
func NewUserRole(userID, roleID string) *UserRole {
	now := time.Now().UTC()
	return &UserRole{
		ID:        uuid.New().String(),
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           string    `json:"id" db:"id"`
	RoleID       string    `json:"role_id" db:"role_id"`
	PermissionID string    `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// NewRolePermission creates a new role permission with default values
func NewRolePermission(roleID, permissionID string) *RolePermission {
	now := time.Now().UTC()
	return &RolePermission{
		ID:           uuid.New().String(),
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Session represents a user's authenticated session
type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Token        string    `json:"token" db:"token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	IP           string    `json:"ip" db:"ip"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// NewSession creates a new session with default values
func NewSession(userID, token, refreshToken string, expiresAt time.Time, ip, userAgent string) *Session {
	now := time.Now().UTC()
	return &Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		IP:           ip,
		UserAgent:    userAgent,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// MFADevice represents a multi-factor authentication device
type MFADevice struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
	Secret    string    `json:"-" db:"secret"`
	Verified  bool      `json:"verified" db:"verified"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewMFADevice creates a new MFA device with default values
func NewMFADevice(userID, deviceType, secret string) *MFADevice {
	now := time.Now().UTC()
	return &MFADevice{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      deviceType,
		Secret:    secret,
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewPasswordReset creates a new password reset with default values
func NewPasswordReset(userID, token string, expiresAt time.Time) *PasswordReset {
	now := time.Now().UTC()
	return &PasswordReset{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}