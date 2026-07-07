SELECT r.name
FROM user_role ur
JOIN role r ON r.id = ur.role_id
WHERE ur.user_id = $1