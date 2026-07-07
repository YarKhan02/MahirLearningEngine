INSERT INTO user_role (user_id, role_id)
SELECT $1, id
FROM role
WHERE name = $2