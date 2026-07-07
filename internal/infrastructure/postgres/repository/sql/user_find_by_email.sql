SELECT id, email, password_hash, is_verified, is_banned, failed_attempts, locked_until
FROM users WHERE email = $1