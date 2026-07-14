SELECT id, email, username, password_hash, is_verified, is_banned, failed_attempts, locked_until
FROM users
WHERE email = $1 OR username = $1
LIMIT 1
