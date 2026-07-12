SELECT s.id
FROM users u
JOIN students s ON s.email = u.email
WHERE u.id = $1
