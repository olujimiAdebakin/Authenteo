-- Rollback all changes in REVERSE order of creation

-- ===== 1. DROP TRIGGERS =====
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- ===== 2. DROP INDEXES (in any order) =====
-- Password reset indexes
DROP INDEX IF EXISTS idx_password_reset_tokens_expires_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;

-- Two-FA indexes
DROP INDEX IF EXISTS idx_two_fa_codes_user_id;
DROP INDEX IF EXISTS idx_two_fa_codes_expires_at;
DROP INDEX IF EXISTS idx_two_fa_codes_email_code;

-- Users indexes
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;

-- ===== 3. DROP TABLES (respecting foreign key dependencies) =====
-- First: tables with foreign keys (referencing users)
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS two_fa_codes;

-- Then: the main users table
DROP TABLE IF EXISTS users;