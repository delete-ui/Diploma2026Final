
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS username VARCHAR(100),
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS verification_code VARCHAR(10),
    ADD COLUMN IF NOT EXISTS verification_expires_at TIMESTAMP;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);