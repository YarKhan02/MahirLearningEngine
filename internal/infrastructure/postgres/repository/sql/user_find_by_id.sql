SELECT id, email, is_verified, is_banned, failed_attempts, locked_until
FROM users WHERE id = $1