package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresRepository implements the Repository interface for PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	// Create the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Connect to the database
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresRepository{
		db: db,
	}, nil
}

// Close closes the database connection
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// Transaction support

// BeginTx begins a transaction
func (r *PostgresRepository) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Store the transaction in the context
	return context.WithValue(ctx, "tx", tx), nil
}

// CommitTx commits a transaction
func (r *PostgresRepository) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return errors.New("no transaction found in context")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackTx rolls back a transaction
func (r *PostgresRepository) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return errors.New("no transaction found in context")
	}

	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// Helper function to get the query executor (db or transaction)
func (r *PostgresRepository) getQueryExecutor(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if ok {
		return tx
	}
	return r.db
}

// User operations

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, active, 
			email_verified, mfa_enabled, last_login, failed_attempts, 
			locked, lock_expiry, created_at, updated_at
		) VALUES (
			:id, :email, :password_hash, :first_name, :last_name, :active, 
			:email_verified, :mfa_enabled, :last_login, :failed_attempts, 
			:locked, :lock_expiry, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID gets a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT * FROM users WHERE id = $1
	`

	var user model.User
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail gets a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`

	var user model.User
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates a user
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET
			email = :email,
			password_hash = :password_hash,
			first_name = :first_name,
			last_name = :last_name,
			active = :active,
			email_verified = :email_verified,
			mfa_enabled = :mfa_enabled,
			last_login = :last_login,
			failed_attempts = :failed_attempts,
			locked = :locked,
			lock_expiry = :lock_expiry,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user
func (r *PostgresRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Role operations

// CreateRole creates a new role
func (r *PostgresRepository) CreateRole(ctx context.Context, role *model.Role) error {
	query := `
		INSERT INTO roles (
			id, name, description, created_at, updated_at
		) VALUES (
			:id, :name, :description, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, role)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetRoleByID gets a role by ID
func (r *PostgresRepository) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	query := `
		SELECT * FROM roles WHERE id = $1
	`

	var role model.Role
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// GetRoleByName gets a role by name
func (r *PostgresRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	query := `
		SELECT * FROM roles WHERE name = $1
	`

	var role model.Role
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &role, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// UpdateRole updates a role
func (r *PostgresRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	query := `
		UPDATE roles SET
			name = :name,
			description = :description,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// DeleteRole deletes a role
func (r *PostgresRepository) DeleteRole(ctx context.Context, id string) error {
	query := `
		DELETE FROM roles WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// GetRolesByUserID gets roles by user ID
func (r *PostgresRepository) GetRolesByUserID(ctx context.Context, userID string) ([]*model.Role, error) {
	query := `
		SELECT r.* FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	var roles []*model.Role
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &roles, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
}

// Permission operations

// CreatePermission creates a new permission
func (r *PostgresRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	query := `
		INSERT INTO permissions (
			id, name, description, resource, action, created_at, updated_at
		) VALUES (
			:id, :name, :description, :resource, :action, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, permission)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

// GetPermissionByID gets a permission by ID
func (r *PostgresRepository) GetPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	query := `
		SELECT * FROM permissions WHERE id = $1
	`

	var permission model.Permission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &permission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("permission not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

// GetPermissionByName gets a permission by name
func (r *PostgresRepository) GetPermissionByName(ctx context.Context, name string) (*model.Permission, error) {
	query := `
		SELECT * FROM permissions WHERE name = $1
	`

	var permission model.Permission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &permission, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("permission not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

// UpdatePermission updates a permission
func (r *PostgresRepository) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	query := `
		UPDATE permissions SET
			name = :name,
			description = :description,
			resource = :resource,
			action = :action,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, permission)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}

	return nil
}

// DeletePermission deletes a permission
func (r *PostgresRepository) DeletePermission(ctx context.Context, id string) error {
	query := `
		DELETE FROM permissions WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}

	return nil
}

// GetPermissionsByRoleID gets permissions by role ID
func (r *PostgresRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error) {
	query := `
		SELECT p.* FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	var permissions []*model.Permission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &permissions, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return permissions, nil
}

// GetPermissionsByUserID gets permissions by user ID
func (r *PostgresRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*model.Permission, error) {
	query := `
		SELECT DISTINCT p.* FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`

	var permissions []*model.Permission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &permissions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return permissions, nil
}

// UserRole operations

// AssignRoleToUser assigns a role to a user
func (r *PostgresRepository) AssignRoleToUser(ctx context.Context, userRole *model.UserRole) error {
	query := `
		INSERT INTO user_roles (
			id, user_id, role_id, created_at, updated_at
		) VALUES (
			:id, :user_id, :role_id, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, userRole)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *PostgresRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	query := `
		DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user role not found")
	}

	return nil
}

// GetUserRolesByUserID gets user roles by user ID
func (r *PostgresRepository) GetUserRolesByUserID(ctx context.Context, userID string) ([]*model.UserRole, error) {
	query := `
		SELECT * FROM user_roles WHERE user_id = $1
	`

	var userRoles []*model.UserRole
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &userRoles, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return userRoles, nil
}

// GetUserRolesByRoleID gets user roles by role ID
func (r *PostgresRepository) GetUserRolesByRoleID(ctx context.Context, roleID string) ([]*model.UserRole, error) {
	query := `
		SELECT * FROM user_roles WHERE role_id = $1
	`

	var userRoles []*model.UserRole
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &userRoles, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return userRoles, nil
}

// RolePermission operations

// AssignPermissionToRole assigns a permission to a role
func (r *PostgresRepository) AssignPermissionToRole(ctx context.Context, rolePermission *model.RolePermission) error {
	query := `
		INSERT INTO role_permissions (
			id, role_id, permission_id, created_at, updated_at
		) VALUES (
			:id, :role_id, :permission_id, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, rolePermission)
	if err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (r *PostgresRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	query := `
		DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role permission not found")
	}

	return nil
}

// GetRolePermissionsByRoleID gets role permissions by role ID
func (r *PostgresRepository) GetRolePermissionsByRoleID(ctx context.Context, roleID string) ([]*model.RolePermission, error) {
	query := `
		SELECT * FROM role_permissions WHERE role_id = $1
	`

	var rolePermissions []*model.RolePermission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &rolePermissions, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return rolePermissions, nil
}

// GetRolePermissionsByPermissionID gets role permissions by permission ID
func (r *PostgresRepository) GetRolePermissionsByPermissionID(ctx context.Context, permissionID string) ([]*model.RolePermission, error) {
	query := `
		SELECT * FROM role_permissions WHERE permission_id = $1
	`

	var rolePermissions []*model.RolePermission
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &rolePermissions, query, permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return rolePermissions, nil
}

// Session operations

// CreateSession creates a new session
func (r *PostgresRepository) CreateSession(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, token, refresh_token, expires_at, ip, user_agent, created_at, updated_at
		) VALUES (
			:id, :user_id, :token, :refresh_token, :expires_at, :ip, :user_agent, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByID gets a session by ID
func (r *PostgresRepository) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	query := `
		SELECT * FROM sessions WHERE id = $1
	`

	var session model.Session
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// GetSessionByToken gets a session by token
func (r *PostgresRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	query := `
		SELECT * FROM sessions WHERE token = $1
	`

	var session model.Session
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &session, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// UpdateSession updates a session
func (r *PostgresRepository) UpdateSession(ctx context.Context, session *model.Session) error {
	query := `
		UPDATE sessions SET
			user_id = :user_id,
			token = :token,
			refresh_token = :refresh_token,
			expires_at = :expires_at,
			ip = :ip,
			user_agent = :user_agent,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, session)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// DeleteSession deletes a session
func (r *PostgresRepository) DeleteSession(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// DeleteSessionsByUserID deletes sessions by user ID
func (r *PostgresRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	query := `
		DELETE FROM sessions WHERE user_id = $1
	`

	exec := r.getQueryExecutor(ctx)
	_, err := exec.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	return nil
}

// MFADevice operations

// CreateMFADevice creates a new MFA device
func (r *PostgresRepository) CreateMFADevice(ctx context.Context, device *model.MFADevice) error {
	query := `
		INSERT INTO mfa_devices (
			id, user_id, type, secret, verified, created_at, updated_at
		) VALUES (
			:id, :user_id, :type, :secret, :verified, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, device)
	if err != nil {
		return fmt.Errorf("failed to create MFA device: %w", err)
	}

	return nil
}

// GetMFADeviceByID gets an MFA device by ID
func (r *PostgresRepository) GetMFADeviceByID(ctx context.Context, id string) (*model.MFADevice, error) {
	query := `
		SELECT * FROM mfa_devices WHERE id = $1
	`

	var device model.MFADevice
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &device, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("MFA device not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get MFA device: %w", err)
	}

	return &device, nil
}

// GetMFADevicesByUserID gets MFA devices by user ID
func (r *PostgresRepository) GetMFADevicesByUserID(ctx context.Context, userID string) ([]*model.MFADevice, error) {
	query := `
		SELECT * FROM mfa_devices WHERE user_id = $1
	`

	var devices []*model.MFADevice
	exec := r.getQueryExecutor(ctx)
	err := sqlx.SelectContext(ctx, exec, &devices, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA devices: %w", err)
	}

	return devices, nil
}

// UpdateMFADevice updates an MFA device
func (r *PostgresRepository) UpdateMFADevice(ctx context.Context, device *model.MFADevice) error {
	query := `
		UPDATE mfa_devices SET
			user_id = :user_id,
			type = :type,
			secret = :secret,
			verified = :verified,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, device)
	if err != nil {
		return fmt.Errorf("failed to update MFA device: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("MFA device not found")
	}

	return nil
}

// DeleteMFADevice deletes an MFA device
func (r *PostgresRepository) DeleteMFADevice(ctx context.Context, id string) error {
	query := `
		DELETE FROM mfa_devices WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete MFA device: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("MFA device not found")
	}

	return nil
}

// PasswordReset operations

// CreatePasswordReset creates a new password reset
func (r *PostgresRepository) CreatePasswordReset(ctx context.Context, reset *model.PasswordReset) error {
	query := `
		INSERT INTO password_resets (
			id, user_id, token, expires_at, used, created_at, updated_at
		) VALUES (
			:id, :user_id, :token, :expires_at, :used, :created_at, :updated_at
		)
	`

	exec := r.getQueryExecutor(ctx)
	_, err := sqlx.NamedExecContext(ctx, exec, query, reset)
	if err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	return nil
}

// GetPasswordResetByToken gets a password reset by token
func (r *PostgresRepository) GetPasswordResetByToken(ctx context.Context, token string) (*model.PasswordReset, error) {
	query := `
		SELECT * FROM password_resets WHERE token = $1
	`

	var reset model.PasswordReset
	exec := r.getQueryExecutor(ctx)
	err := sqlx.GetContext(ctx, exec, &reset, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("password reset not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get password reset: %w", err)
	}

	return &reset, nil
}

// UpdatePasswordReset updates a password reset
func (r *PostgresRepository) UpdatePasswordReset(ctx context.Context, reset *model.PasswordReset) error {
	query := `
		UPDATE password_resets SET
			user_id = :user_id,
			token = :token,
			expires_at = :expires_at,
			used = :used,
			updated_at = :updated_at
		WHERE id = :id
	`

	exec := r.getQueryExecutor(ctx)
	result, err := sqlx.NamedExecContext(ctx, exec, query, reset)
	if err != nil {
		return fmt.Errorf("failed to update password reset: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("password reset not found")
	}

	return nil
}

// DeletePasswordReset deletes a password reset
func (r *PostgresRepository) DeletePasswordReset(ctx context.Context, id string) error {
	query := `
		DELETE FROM password_resets WHERE id = $1
	`

	exec := r.getQueryExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete password reset: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("password reset not found")
	}

	return nil
}

// DeleteExpiredPasswordResets deletes expired password resets
func (r *PostgresRepository) DeleteExpiredPasswordResets(ctx context.Context) error {
	query := `
		DELETE FROM password_resets WHERE expires_at < NOW() OR used = true
	`

	exec := r.getQueryExecutor(ctx)
	_, err := exec.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired password resets: %w", err)
	}

	return nil
}
