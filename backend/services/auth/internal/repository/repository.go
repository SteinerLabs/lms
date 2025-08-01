package repository

import (
	"context"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/model"
)

// Repository defines the interface for database operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error

	// Role operations
	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByID(ctx context.Context, id string) (*model.Role, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id string) error
	GetRolesByUserID(ctx context.Context, userID string) ([]*model.Role, error)

	// Permission operations
	CreatePermission(ctx context.Context, permission *model.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*model.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*model.Permission, error)
	UpdatePermission(ctx context.Context, permission *model.Permission) error
	DeletePermission(ctx context.Context, id string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID string) ([]*model.Permission, error)

	// UserRole operations
	AssignRoleToUser(ctx context.Context, userRole *model.UserRole) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRolesByUserID(ctx context.Context, userID string) ([]*model.UserRole, error)
	GetUserRolesByRoleID(ctx context.Context, roleID string) ([]*model.UserRole, error)

	// RolePermission operations
	AssignPermissionToRole(ctx context.Context, rolePermission *model.RolePermission) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	GetRolePermissionsByRoleID(ctx context.Context, roleID string) ([]*model.RolePermission, error)
	GetRolePermissionsByPermissionID(ctx context.Context, permissionID string) ([]*model.RolePermission, error)

	// Session operations
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
	DeleteSession(ctx context.Context, id string) error
	DeleteSessionsByUserID(ctx context.Context, userID string) error

	// MFADevice operations
	CreateMFADevice(ctx context.Context, device *model.MFADevice) error
	GetMFADeviceByID(ctx context.Context, id string) (*model.MFADevice, error)
	GetMFADevicesByUserID(ctx context.Context, userID string) ([]*model.MFADevice, error)
	UpdateMFADevice(ctx context.Context, device *model.MFADevice) error
	DeleteMFADevice(ctx context.Context, id string) error

	// PasswordReset operations
	CreatePasswordReset(ctx context.Context, reset *model.PasswordReset) error
	GetPasswordResetByToken(ctx context.Context, token string) (*model.PasswordReset, error)
	UpdatePasswordReset(ctx context.Context, reset *model.PasswordReset) error
	DeletePasswordReset(ctx context.Context, id string) error
	DeleteExpiredPasswordResets(ctx context.Context) error

	// Transaction support
	BeginTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error

	// Close the repository
	Close() error
}
