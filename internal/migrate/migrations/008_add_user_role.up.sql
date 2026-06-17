
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';

ALTER TABLE users
    ADD CONSTRAINT check_role_valid
        CHECK (role IN ('user', 'admin', 'developer'));

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

UPDATE users SET role = 'user' WHERE role IS NULL;