SELECT s.id
FROM users u
JOIN students s ON s.username = u.username
WHERE u.id = $1
