-- Table: users (global scope)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(100),
    phone_number VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- 'user', 'superadmin'
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Table: user_metadata (optional, JSONB flexible info)
CREATE TABLE IF NOT EXISTS user_metadata (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Index for faster lookup (optional)
CREATE INDEX IF NOT EXISTS idx_user_metadata_user_id ON user_metadata(user_id);