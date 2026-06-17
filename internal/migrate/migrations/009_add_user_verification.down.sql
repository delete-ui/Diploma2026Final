DROP INDEX IF EXISTS idx_users_username;
ALTER TABLE users
DROP COLUMN IF EXISTS username,
    DROP COLUMN IF EXISTS is_verified,
    DROP COLUMN IF EXISTS verification_code,
    DROP COLUMN IF EXISTS verification_expires_at;