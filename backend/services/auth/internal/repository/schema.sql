-- Auth Service Database Schema

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMP,
    failed_attempts INT NOT NULL DEFAULT 0,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    lock_expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (resource, action)
);

-- Create user_roles table (many-to-many relationship between users and roles)
CREATE TABLE IF NOT EXISTS user_roles (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

-- Create role_permissions table (many-to-many relationship between roles and permissions)
CREATE TABLE IF NOT EXISTS role_permissions (
    id VARCHAR(36) PRIMARY KEY,
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    ip VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create mfa_devices table
CREATE TABLE IF NOT EXISTS mfa_devices (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create password_resets table
CREATE TABLE IF NOT EXISTS password_resets (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles (name);
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions (name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions (resource, action);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions (permission_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX IF NOT EXISTS idx_mfa_devices_user_id ON mfa_devices (user_id);
CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON password_resets (user_id);
CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets (token);

-- Create default roles
INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'admin', 'Administrator with full access', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'user', 'Regular user with limited access', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000003', 'instructor', 'Instructor with course management access', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Create default permissions
INSERT INTO permissions (id, name, description, resource, action, created_at, updated_at)
VALUES 
    -- User permissions
    ('00000000-0000-0000-0000-000000000001', 'user:create', 'Create users', 'user', 'create', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'user:read', 'Read users', 'user', 'read', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000003', 'user:update', 'Update users', 'user', 'update', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000004', 'user:delete', 'Delete users', 'user', 'delete', NOW(), NOW()),
    
    -- Role permissions
    ('00000000-0000-0000-0000-000000000005', 'role:create', 'Create roles', 'role', 'create', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000006', 'role:read', 'Read roles', 'role', 'read', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000007', 'role:update', 'Update roles', 'role', 'update', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000008', 'role:delete', 'Delete roles', 'role', 'delete', NOW(), NOW()),
    
    -- Permission permissions
    ('00000000-0000-0000-0000-000000000009', 'permission:create', 'Create permissions', 'permission', 'create', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000010', 'permission:read', 'Read permissions', 'permission', 'read', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000011', 'permission:update', 'Update permissions', 'permission', 'update', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000012', 'permission:delete', 'Delete permissions', 'permission', 'delete', NOW(), NOW()),
    
    -- Course permissions
    ('00000000-0000-0000-0000-000000000013', 'course:create', 'Create courses', 'course', 'create', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000014', 'course:read', 'Read courses', 'course', 'read', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000015', 'course:update', 'Update courses', 'course', 'update', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000016', 'course:delete', 'Delete courses', 'course', 'delete', NOW(), NOW())
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to roles
INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
VALUES 
    -- Admin role has all permissions
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000003', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000005', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000006', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000007', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000008', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000009', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000010', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000011', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000012', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000013', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000014', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000015', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000016', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000016', NOW(), NOW()),
    
    -- User role has read permissions only
    ('00000000-0000-0000-0000-000000000017', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000018', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000006', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000019', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000010', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000020', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000014', NOW(), NOW()),
    
    -- Instructor role has course management permissions
    ('00000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000002', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000022', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000006', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000023', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000010', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000024', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000013', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000025', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000014', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000026', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000015', NOW(), NOW())
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Create a default admin user (password: admin123)
INSERT INTO users (
    id, email, password_hash, first_name, last_name, 
    active, email_verified, mfa_enabled, created_at, updated_at
)
VALUES (
    '00000000-0000-0000-0000-000000000001', 
    'admin@example.com', 
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- bcrypt hash for 'admin123'
    'Admin',
    'User',
    TRUE,
    TRUE,
    FALSE,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Assign admin role to admin user
INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    NOW(),
    NOW()
)
ON CONFLICT (user_id, role_id) DO NOTHING;