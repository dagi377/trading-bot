-- Add users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Insert sample users
-- Note: In production, passwords should be properly hashed
INSERT INTO users (username, password_hash, email) VALUES
    ('admin', '$2a$10$X7UrH5Zx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx', 'admin@example.com'),
    ('demo', '$2a$10$X7UrH5Zx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx', 'demo@example.com'),
    ('test', '$2a$10$X7UrH5Zx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx.5Yx', 'test@example.com')
ON CONFLICT (username) DO NOTHING; 
