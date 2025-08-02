package service

import (
	"context"
	"fmt"
	"github.com/SteinerLabs/lms/backend/services/auth/proto/gen/proto"
	"time"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/model"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/repository"
	"github.com/SteinerLabs/lms/backend/shared/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthService implements the AuthService gRPC service
type AuthService struct {
	proto.UnimplementedAuthServiceServer
	repo   repository.Repository
	config *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config) (*AuthService, error) {
	// Create a new repository
	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &AuthService{
		repo:   repo,
		config: cfg,
	}, nil
}

// Close closes the service
func (s *AuthService) Close() error {
	return s.repo.Close()
}

// User management

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Create a new user
	user := model.NewUser(req.Email, string(hashedPassword), req.FirstName, req.LastName)

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to begin transaction")
	}

	// Create the user
	err = s.repo.CreateUser(txCtx, user)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// Assign default role (user)
	defaultRoleID := "00000000-0000-0000-0000-000000000002" // user role
	userRole := model.NewUserRole(user.ID, defaultRoleID)
	err = s.repo.AssignRoleToUser(txCtx, userRole)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to assign default role")
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to commit transaction")
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
	// TODO: Publish event to Kafka

	// Convert to proto user
	protoUser := &auth.User{
		Id:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Active:        user.Active,
		EmailVerified: user.EmailVerified,
		MfaEnabled:    user.MFAEnabled,
		LastLogin:     timestamppb.New(user.LastLogin),
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
	}

	return protoUser, nil
}

// GetUser gets a user by ID
func (s *AuthService) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.User, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Convert to proto user
	protoUser := &auth.User{
		Id:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Active:        user.Active,
		EmailVerified: user.EmailVerified,
		MfaEnabled:    user.MFAEnabled,
		LastLogin:     timestamppb.New(user.LastLogin),
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
	}

	return protoUser, nil
}

// UpdateUser updates a user
func (s *AuthService) UpdateUser(ctx context.Context, req *auth.UpdateUserRequest) (*auth.User, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Update the user
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	user.Active = req.Active
	user.UpdatedAt = time.Now().UTC()

	// Update the user
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
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
	// TODO: Publish event to Kafka

	// Convert to proto user
	protoUser := &auth.User{
		Id:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Active:        user.Active,
		EmailVerified: user.EmailVerified,
		MfaEnabled:    user.MFAEnabled,
		LastLogin:     timestamppb.New(user.LastLogin),
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
	}

	return protoUser, nil
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, req *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to begin transaction")
	}

	// Delete the user
	err = s.repo.DeleteUser(txCtx, req.Id)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	// Publish user deleted event
	userDeletedEvent := &events.UserDeletedEvent{
		ID:        user.ID,
		DeletedAt: time.Now().UTC(),
	}
	event := events.NewEvent("user.deleted", "auth-service", userDeletedEvent, "", "")
	// TODO: Publish event to Kafka

	return &auth.DeleteUserResponse{
		Success: true,
	}, nil
}

// Authentication

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Parse the token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(s.config.JWT.Secret), nil
	})

	// Check for parsing errors
	if err != nil {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Check if the token is valid
	if !token.Valid {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Get the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Get the user ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Get the user's permissions
	permissions, err := s.repo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get permissions")
	}

	// Convert permissions to strings
	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = permission.Name
	}

	return &auth.ValidateTokenResponse{
		Valid:       true,
		UserId:      userID,
		Permissions: permissionStrings,
	}, nil
}

// GetUserPermissions gets a user's permissions
func (s *AuthService) GetUserPermissions(ctx context.Context, req *auth.GetUserPermissionsRequest) (*auth.GetUserPermissionsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Get the user's permissions
	permissions, err := s.repo.GetPermissionsByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get permissions")
	}

	// Convert permissions to strings
	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = permission.Name
	}

	return &auth.GetUserPermissionsResponse{
		Permissions: permissionStrings,
	}, nil
}

// Role management

// CreateRole creates a new role
func (s *AuthService) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Check if role already exists
	existingRole, err := s.repo.GetRoleByName(ctx, req.Name)
	if err == nil && existingRole != nil {
		return nil, status.Error(codes.AlreadyExists, "role already exists")
	}

	// Create a new role
	role := model.NewRole(req.Name, req.Description)

	// Create the role
	err = s.repo.CreateRole(ctx, role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create role")
	}

	// Convert to proto role
	protoRole := &auth.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   timestamppb.New(role.UpdatedAt),
	}

	return protoRole, nil
}

// GetRole gets a role by ID
func (s *AuthService) GetRole(ctx context.Context, req *auth.GetRoleRequest) (*auth.Role, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the role
	role, err := s.repo.GetRoleByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	// Convert to proto role
	protoRole := &auth.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   timestamppb.New(role.UpdatedAt),
	}

	return protoRole, nil
}

// UpdateRole updates a role
func (s *AuthService) UpdateRole(ctx context.Context, req *auth.UpdateRoleRequest) (*auth.Role, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the role
	role, err := s.repo.GetRoleByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	// Update the role
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.UpdatedAt = time.Now().UTC()

	// Update the role
	err = s.repo.UpdateRole(ctx, role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update role")
	}

	// Convert to proto role
	protoRole := &auth.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   timestamppb.New(role.UpdatedAt),
	}

	return protoRole, nil
}

// DeleteRole deletes a role
func (s *AuthService) DeleteRole(ctx context.Context, req *auth.DeleteRoleRequest) (*auth.DeleteRoleResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get the role
	_, err := s.repo.GetRoleByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	// Begin transaction
	txCtx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to begin transaction")
	}

	// Delete the role
	err = s.repo.DeleteRole(txCtx, req.Id)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to delete role")
	}

	// Commit transaction
	err = s.repo.CommitTx(txCtx)
	if err != nil {
		s.repo.RollbackTx(txCtx)
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	return &auth.DeleteRoleResponse{
		Success: true,
	}, nil
}

// AssignRoleToUser assigns a role to a user
func (s *AuthService) AssignRoleToUser(ctx context.Context, req *auth.AssignRoleToUserRequest) (*auth.AssignRoleToUserResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "role_id is required")
	}

	// Get the user
	_, err := s.repo.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Get the role
	_, err = s.repo.GetRoleByID(ctx, req.RoleId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	// Create a new user role
	userRole := model.NewUserRole(req.UserId, req.RoleId)

	// Assign the role to the user
	err = s.repo.AssignRoleToUser(ctx, userRole)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to assign role to user")
	}

	return &auth.AssignRoleToUserResponse{
		Success: true,
	}, nil
}

// RemoveRoleFromUser removes a role from a user
func (s *AuthService) RemoveRoleFromUser(ctx context.Context, req *auth.RemoveRoleFromUserRequest) (*auth.RemoveRoleFromUserResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "role_id is required")
	}

	// Get the user
	_, err := s.repo.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Get the role
	_, err = s.repo.GetRoleByID(ctx, req.RoleId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "role not found")
	}

	// Remove the role from the user
	err = s.repo.RemoveRoleFromUser(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to remove role from user")
	}

	return &auth.RemoveRoleFromUserResponse{
		Success: true,
	}, nil
}

// Helper functions

// generateToken generates a JWT token
func (s *AuthService) generateToken(userID string, expiresIn time.Duration) (string, error) {
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
func (s *AuthService) generateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, time.Duration(s.config.JWT.AccessTokenTTL)*time.Minute)
}

// generateRefreshToken generates a refresh token
func (s *AuthService) generateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, time.Duration(s.config.JWT.RefreshTokenTTL)*24*time.Hour)
}
