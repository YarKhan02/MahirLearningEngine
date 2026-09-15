SELECT COALESCE(NULLIF(st.full_name, ''), u.email)
FROM users u
LEFT JOIN students st ON st.username = u.username
WHERE u.id = $1
