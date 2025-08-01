package service

import (
	"context"
	"testing"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/event"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/model"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/repository"
)

// MockRepository is a mock implementation of the Repository interface for testing
type MockRepository struct {
	users           map[string]*model.User
	usersByEmail    map[string]*model.User
	roles           map[string]*model.Role
	userRoles       map[string]map[string]bool
	sessions        map[string]*model.Session
	sessionsByToken map[string]*model.Session
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:           make(map[string]*model.User),
		usersByEmail:    make(map[string]*model.User),
		roles:           make(map[string]*model.Role),
		userRoles:       make(map[string]map[string]bool),
		sessions:        make(map[string]*model.Session),
		sessionsByToken: make(map[string]*model.Session),
	}
}

// Implement the Repository interface methods for the mock repository
func (r *MockRepository) CreateUser(ctx context.Context, user *model.User) error {
	r.users[user.ID] = user
	r.usersByEmail[user.Email] = user
	return nil
}

func (r *MockRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (r *MockRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (r *MockRepository) UpdateUser(ctx context.Context, user *model.User) error {
	r.users[user.ID] = user
	r.usersByEmail[user.Email] = user
	return nil
}

func (r *MockRepository) DeleteUser(ctx context.Context, id string) error {
	user, ok := r.users[id]
	if !ok {
		return repository.ErrNotFound
	}
	delete(r.usersByEmail, user.Email)
	delete(r.users, id)
	return nil
}

func (r *MockRepository) CreateRole(ctx context.Context, role *model.Role) error {
	r.roles[role.ID] = role
	return nil
}

func (r *MockRepository) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	role, ok := r.roles[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return role, nil
}

func (r *MockRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	for _, role := range r.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *MockRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	r.roles[role.ID] = role
	return nil
}

func (r *MockRepository) DeleteRole(ctx context.Context, id string) error {
	_, ok := r.roles[id]
	if !ok {
		return repository.ErrNotFound
	}
	delete(r.roles, id)
	return nil
}

func (r *MockRepository) GetRolesByUserID(ctx context.Context, userID string) ([]*model.Role, error) {
	roleIDs, ok := r.userRoles[userID]
	if !ok {
		return []*model.Role{}, nil
	}

	roles := make([]*model.Role, 0, len(roleIDs))
	for roleID := range roleIDs {
		role, ok := r.roles[roleID]
		if ok {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func (r *MockRepository) AssignRoleToUser(ctx context.Context, userRole *model.UserRole) error {
	_, ok := r.userRoles[userRole.UserID]
	if !ok {
		r.userRoles[userRole.UserID] = make(map[string]bool)
	}
	r.userRoles[userRole.UserID][userRole.RoleID] = true
	return nil
}

func (r *MockRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	roleIDs, ok := r.userRoles[userID]
	if !ok {
		return repository.ErrNotFound
	}
	delete(roleIDs, roleID)
	return nil
}

func (r *MockRepository) GetUserRolesByUserID(ctx context.Context, userID string) ([]*model.UserRole, error) {
	roleIDs, ok := r.userRoles[userID]
	if !ok {
		return []*model.UserRole{}, nil
	}

	userRoles := make([]*model.UserRole, 0, len(roleIDs))
	for roleID := range roleIDs {
		userRoles = append(userRoles, &model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return userRoles, nil
}

func (r *MockRepository) GetUserRolesByRoleID(ctx context.Context, roleID string) ([]*model.UserRole, error) {
	userRoles := make([]*model.UserRole, 0)
	for userID, roleIDs := range r.userRoles {
		if roleIDs[roleID] {
			userRoles = append(userRoles, &model.UserRole{
				UserID: userID,
				RoleID: roleID,
			})
		}
	}
	return userRoles, nil
}

func (r *MockRepository) CreateSession(ctx context.Context, session *model.Session) error {
	r.sessions[session.ID] = session
	r.sessionsByToken[session.Token] = session
	r.sessionsByToken[session.RefreshToken] = session
	return nil
}

func (r *MockRepository) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	session, ok := r.sessions[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return session, nil
}

func (r *MockRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	session, ok := r.sessionsByToken[token]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return session, nil
}

func (r *MockRepository) UpdateSession(ctx context.Context, session *model.Session) error {
	r.sessions[session.ID] = session
	r.sessionsByToken[session.Token] = session
	r.sessionsByToken[session.RefreshToken] = session
	return nil
}

func (r *MockRepository) DeleteSession(ctx context.Context, id string) error {
	session, ok := r.sessions[id]
	if !ok {
		return repository.ErrNotFound
	}
	delete(r.sessionsByToken, session.Token)
	delete(r.sessionsByToken, session.RefreshToken)
	delete(r.sessions, id)
	return nil
}

func (r *MockRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	for id, session := range r.sessions {
		if session.UserID == userID {
			delete(r.sessionsByToken, session.Token)
			delete(r.sessionsByToken, session.RefreshToken)
			delete(r.sessions, id)
		}
	}
	return nil
}

// Stub implementations for the remaining repository methods
func (r *MockRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return nil
}

func (r *MockRepository) GetPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	return nil, nil
}

func (r *MockRepository) GetPermissionByName(ctx context.Context, name string) (*model.Permission, error) {
	return nil, nil
}

func (r *MockRepository) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	return nil
}

func (r *MockRepository) DeletePermission(ctx context.Context, id string) error {
	return nil
}

func (r *MockRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error) {
	return nil, nil
}

func (r *MockRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*model.Permission, error) {
	return []*model.Permission{}, nil
}

func (r *MockRepository) AssignPermissionToRole(ctx context.Context, rolePermission *model.RolePermission) error {
	return nil
}

func (r *MockRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (r *MockRepository) GetRolePermissionsByRoleID(ctx context.Context, roleID string) ([]*model.RolePermission, error) {
	return nil, nil
}

func (r *MockRepository) GetRolePermissionsByPermissionID(ctx context.Context, permissionID string) ([]*model.RolePermission, error) {
	return nil, nil
}

func (r *MockRepository) CreateMFADevice(ctx context.Context, device *model.MFADevice) error {
	return nil
}

func (r *MockRepository) GetMFADeviceByID(ctx context.Context, id string) (*model.MFADevice, error) {
	return nil, nil
}

func (r *MockRepository) GetMFADevicesByUserID(ctx context.Context, userID string) ([]*model.MFADevice, error) {
	return nil, nil
}

func (r *MockRepository) UpdateMFADevice(ctx context.Context, device *model.MFADevice) error {
	return nil
}

func (r *MockRepository) DeleteMFADevice(ctx context.Context, id string) error {
	return nil
}

func (r *MockRepository) CreatePasswordReset(ctx context.Context, reset *model.PasswordReset) error {
	return nil
}

func (r *MockRepository) GetPasswordResetByToken(ctx context.Context, token string) (*model.PasswordReset, error) {
	return nil, nil
}

func (r *MockRepository) UpdatePasswordReset(ctx context.Context, reset *model.PasswordReset) error {
	return nil
}

func (r *MockRepository) DeletePasswordReset(ctx context.Context, id string) error {
	return nil
}

func (r *MockRepository) DeleteExpiredPasswordResets(ctx context.Context) error {
	return nil
}

func (r *MockRepository) BeginTx(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (r *MockRepository) CommitTx(ctx context.Context) error {
	return nil
}

func (r *MockRepository) RollbackTx(ctx context.Context) error {
	return nil
}

func (r *MockRepository) Close() error {
	return nil
}

// TestAuthService tests the auth service
func TestAuthService(t *testing.T) {
	// Create a mock repository
	repo := NewMockRepository()

	// Create a mock event publisher
	publisher := event.NewMockPublisher()

	// Create a test config
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			AccessTokenTTL:  15,
			RefreshTokenTTL: 7,
			Issuer:          "test-issuer",
			Audience:        "test-audience",
		},
	}

	// Create the auth service with the mock repository and publisher
	service := &AuthServiceImpl{
		repo:      repo,
		config:    cfg,
		publisher: publisher,
	}

	// Test user creation
	t.Run("CreateUser", func(t *testing.T) {
		// Create a new user
		user, err := service.CreateUser(context.Background(), "test@example.com", "password123", "Test", "User")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Check that the user was created correctly
		if user.Email != "test@example.com" {
			t.Errorf("Expected email to be test@example.com, got %s", user.Email)
		}
		if user.FirstName != "Test" {
			t.Errorf("Expected first name to be Test, got %s", user.FirstName)
		}
		if user.LastName != "User" {
			t.Errorf("Expected last name to be User, got %s", user.LastName)
		}
		if !user.Active {
			t.Errorf("Expected user to be active")
		}
		if user.EmailVerified {
			t.Errorf("Expected user to not be email verified")
		}
		if user.MFAEnabled {
			t.Errorf("Expected user to not have MFA enabled")
		}

		// Check that the user was stored in the repository
		storedUser, err := repo.GetUserByEmail(context.Background(), "test@example.com")
		if err != nil {
			t.Fatalf("Failed to get user from repository: %v", err)
		}
		if storedUser.ID != user.ID {
			t.Errorf("Expected user ID to be %s, got %s", user.ID, storedUser.ID)
		}

		// Check that a user created event was published
		events := publisher.GetEventsByType("user.created")
		if len(events) != 1 {
			t.Fatalf("Expected 1 user.created event, got %d", len(events))
		}
	})

	// Test user login
	t.Run("Login", func(t *testing.T) {
		// Login the user
		session, err := service.Login(context.Background(), "test@example.com", "password123", "127.0.0.1", "test-user-agent")
		if err != nil {
			t.Fatalf("Failed to login user: %v", err)
		}

		// Check that the session was created correctly
		if session.UserID == "" {
			t.Errorf("Expected session to have a user ID")
		}
		if session.Token == "" {
			t.Errorf("Expected session to have a token")
		}
		if session.RefreshToken == "" {
			t.Errorf("Expected session to have a refresh token")
		}
		if session.ExpiresAt.IsZero() {
			t.Errorf("Expected session to have an expiration time")
		}
		if session.IP != "127.0.0.1" {
			t.Errorf("Expected session IP to be 127.0.0.1, got %s", session.IP)
		}
		if session.UserAgent != "test-user-agent" {
			t.Errorf("Expected session user agent to be test-user-agent, got %s", session.UserAgent)
		}

		// Check that the session was stored in the repository
		storedSession, err := repo.GetSessionByID(context.Background(), session.ID)
		if err != nil {
			t.Fatalf("Failed to get session from repository: %v", err)
		}
		if storedSession.ID != session.ID {
			t.Errorf("Expected session ID to be %s, got %s", session.ID, storedSession.ID)
		}

		// Check that a user logged in event was published
		events := publisher.GetEventsByType("user.login")
		if len(events) != 1 {
			t.Fatalf("Expected 1 user.login event, got %d", len(events))
		}

		// Test token validation
		userID, permissions, err := service.ValidateToken(context.Background(), session.Token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}
		if userID == "" {
			t.Errorf("Expected user ID to be returned from token validation")
		}
		if len(permissions) != 0 {
			t.Errorf("Expected 0 permissions, got %d", len(permissions))
		}

		// Test logout
		err = service.Logout(context.Background(), session.Token)
		if err != nil {
			t.Fatalf("Failed to logout user: %v", err)
		}

		// Check that the session was removed from the repository
		_, err = repo.GetSessionByID(context.Background(), session.ID)
		if err == nil {
			t.Errorf("Expected session to be removed from repository")
		}

		// Check that a user logged out event was published
		events = publisher.GetEventsByType("user.logout")
		if len(events) != 1 {
			t.Fatalf("Expected 1 user.logout event, got %d", len(events))
		}
	})
}
