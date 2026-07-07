SELECT
    id,
    title,
    description,
    level,
    duration,
    is_active
FROM course
ORDER BY created_at DESC;