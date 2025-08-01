package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/event"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/model"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/repository"
	"github.com/SteinerLabs/lms/backend/shared/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl implements the core functionality of the auth service
type AuthServiceImpl struct {
	repo      repository.Repository
	config    *config.Config
	publisher event.Publisher
}

// NewAuthServiceImpl creates a new AuthServiceImpl
func NewAuthServiceImpl(cfg *config.Config) (*AuthServiceImpl, error) {
	// Create a new repository
	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Create a new event publisher
	publisher, err := event.NewKafkaPublisher(cfg)
	if err != nil {
		repo.Close()
		return nil, fmt.Errorf("failed to create event publisher: %w", err)
	}

	return &AuthServiceImpl{
		repo:      repo,
		config:    cfg,
		publisher: publisher,
	}, nil
}

// Close closes the service
func (s *AuthServiceImpl) Close() error {
	// Close the repository
	repoErr := s.repo.Close()

	// Close the publisher
	pubErr := s.publisher.Close()

	// Return the first error encountered
	if repoErr != nil {
		return repoErr
	}
	return pubErr
}

// User management

// CreateUser creates a new user
func (s *AuthServiceImpl) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*model.User, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new user
	user := model.NewUser(email, string(hashedPassword), firstName, lastName)

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create the user
	err = s.repo.CreateUser(txCtx, user)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role (user)
	defaultRoleID := "00000000-0000-0000-0000-000000000002" // user role
	userRole := model.NewUserRole(user.ID, defaultRoleID)
	err = s.repo.AssignRoleToUser(txCtx, userRole)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish user created event
	userCreatedEvent := &events.UserCreatedEvent{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
	event := events.NewEvent("user.created", "auth-service", userCreatedEvent, "", "")
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to publish user created event: %v\n", err)
	}

	return user, nil
}

// GetUser gets a user by ID
func (s *AuthServiceImpl) GetUser(ctx context.Context, id string) (*model.User, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (s *AuthServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Get the user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user
func (s *AuthServiceImpl) UpdateUser(ctx context.Context, id, email, firstName, lastName string, active bool) (*model.User, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update the user
	if email != "" {
		user.Email = email
	}
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	user.Active = active
	user.UpdatedAt = time.Now().UTC()

	// Update the user
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish user updated event
	userUpdatedEvent := &events.UserUpdatedEvent{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UpdatedAt: user.UpdatedAt,
	}
	event := events.NewEvent("user.updated", "auth-service", userUpdatedEvent, "", "")
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to publish user updated event: %v\n", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *AuthServiceImpl) DeleteUser(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return errors.New("id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Delete the user
	err = s.repo.DeleteUser(txCtx, id)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish user deleted event
	userDeletedEvent := &events.UserDeletedEvent{
		ID:        user.ID,
		DeletedAt: time.Now().UTC(),
	}
	event := events.NewEvent("user.deleted", "auth-service", userDeletedEvent, "", "")
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to publish user deleted event: %v\n", err)
	}

	return nil
}

// Authentication

// Login authenticates a user and returns a session
func (s *AuthServiceImpl) Login(ctx context.Context, email, password, ip, userAgent string) (*model.Session, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Get the user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if the user is active
	if !user.Active {
		return nil, errors.New("user is not active")
	}

	// Check if the user is locked
	if user.Locked && user.LockExpiry.After(time.Now()) {
		return nil, errors.New("account is locked")
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Increment failed attempts
		user.FailedAttempts++

		// Lock the account if too many failed attempts
		if user.FailedAttempts >= 5 {
			user.Locked = true
			user.LockExpiry = time.Now().Add(15 * time.Minute)
		}

		// Update the user
		s.repo.UpdateUser(ctx, user)

		return nil, errors.New("invalid email or password")
	}

	// Reset failed attempts
	user.FailedAttempts = 0
	user.Locked = false
	user.LastLogin = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	// Update the user
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create a new session
	expiresAt := time.Now().Add(time.Duration(s.config.JWT.AccessTokenTTL) * time.Minute)
	session := model.NewSession(user.ID, accessToken, refreshToken, expiresAt, ip, userAgent)

	// Create the session
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Publish user logged in event
	userLoggedInEvent := &events.UserLoggedInEvent{
		ID:        user.ID,
		Email:     user.Email,
		IP:        ip,
		UserAgent: userAgent,
		LoginAt:   time.Now().UTC(),
	}
	event := events.NewEvent("user.login", "auth-service", userLoggedInEvent, "", "")
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to publish user logged in event: %v\n", err)
	}

	return session, nil
}

// Logout logs out a user
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// Validate input
	if token == "" {
		return errors.New("token is required")
	}

	// Get the session
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete the session
	err = s.repo.DeleteSession(ctx, session.ID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Publish user logged out event
	userLoggedOutEvent := &events.UserLoggedOutEvent{
		ID:       user.ID,
		Email:    user.Email,
		LogoutAt: time.Now().UTC(),
	}
	event := events.NewEvent("user.logout", "auth-service", userLoggedOutEvent, "", "")
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to publish user logged out event: %v\n", err)
	}

	return nil
}

// ValidateToken validates a JWT token
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (string, []string, error) {
	// Validate input
	if token == "" {
		return "", nil, errors.New("token is required")
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(s.config.JWT.Secret), nil
	})

	// Check for parsing errors
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !parsedToken.Valid {
		return "", nil, errors.New("invalid token")
	}

	// Get the claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.New("invalid token claims")
	}

	// Get the user ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", nil, errors.New("invalid user ID in token")
	}

	// Get the user's permissions
	permissions, err := s.repo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Convert permissions to strings
	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = permission.Name
	}

	return userID, permissionStrings, nil
}

// RefreshToken refreshes a JWT token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	// Validate input
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	// Get the session
	session, err := s.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if the session is expired
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate a new refresh token
	newRefreshToken, err := s.generateRefreshToken(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update the session
	session.Token = accessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(time.Duration(s.config.JWT.AccessTokenTTL) * time.Minute)
	session.UpdatedAt = time.Now().UTC()

	// Update the session
	err = s.repo.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// GetUserPermissions gets a user's permissions
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user's permissions
	permissions, err := s.repo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Convert permissions to strings
	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = permission.Name
	}

	return permissionStrings, nil
}

// Role management

// CreateRole creates a new role
func (s *AuthServiceImpl) CreateRole(ctx context.Context, name, description string) (*model.Role, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Check if role already exists
	existingRole, err := s.repo.GetRoleByName(ctx, name)
	if err == nil && existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Create a new role
	role := model.NewRole(name, description)

	// Create the role
	err = s.repo.CreateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// GetRole gets a role by ID
func (s *AuthServiceImpl) GetRole(ctx context.Context, id string) (*model.Role, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get the role
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// UpdateRole updates a role
func (s *AuthServiceImpl) UpdateRole(ctx context.Context, id, name, description string) (*model.Role, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get the role
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Update the role
	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}
	role.UpdatedAt = time.Now().UTC()

	// Update the role
	err = s.repo.UpdateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// DeleteRole deletes a role
func (s *AuthServiceImpl) DeleteRole(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return errors.New("id is required")
	}

	// Get the role
	_, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Delete the role
	err = s.repo.DeleteRole(txCtx, id)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AssignRoleToUser assigns a role to a user
func (s *AuthServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// Validate input
	if userID == "" {
		return errors.New("user ID is required")
	}
	if roleID == "" {
		return errors.New("role ID is required")
	}

	// Get the user
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get the role
	_, err = s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Create a new user role
	userRole := model.NewUserRole(userID, roleID)

	// Assign the role to the user
	err = s.repo.AssignRoleToUser(ctx, userRole)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *AuthServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	// Validate input
	if userID == "" {
		return errors.New("user ID is required")
	}
	if roleID == "" {
		return errors.New("role ID is required")
	}

	// Get the user
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get the role
	_, err = s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Remove the role from the user
	err = s.repo.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// Helper functions

// generateToken generates a JWT token
func (s *AuthServiceImpl) generateToken(userID string, expiresIn time.Duration) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.config.JWT.Issuer,
		"aud": s.config.JWT.Audience,
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateAccessToken generates an access token
func (s *AuthServiceImpl) generateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, time.Duration(s.config.JWT.AccessTokenTTL)*time.Minute)
}

// generateRefreshToken generates a refresh token
func (s *AuthServiceImpl) generateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, time.Duration(s.config.JWT.RefreshTokenTTL)*24*time.Hour)
}
