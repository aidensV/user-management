-- =============================================
-- 1. USERS
-- Menyimpan data dasar user
-- =============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(50) UNIQUE,
    department VARCHAR(100),
    position VARCHAR(100),
    phone VARCHAR(20),
    
    -- Role sementara untuk P0 (nanti migrate ke user_roles)
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'supervisor', 'warehouse_operator', 'viewer')),
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_locked BOOLEAN DEFAULT false,
    failed_login_attempts INT DEFAULT 0,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    
    -- Soft delete
    deleted_at TIMESTAMP NULL,
    deleted_by UUID NULL,
    
    -- Audit
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by UUID NULL,
    updated_by UUID NULL
);

-- Indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);