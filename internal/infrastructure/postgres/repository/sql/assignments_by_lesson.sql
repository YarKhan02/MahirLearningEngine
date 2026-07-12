SELECT
    id,
    lesson_id,
    title,
    COALESCE(description, ''),
    COALESCE(starter_code, ''),
    due_date,
    total_marks,
    created_at
FROM assignments
WHERE lesson_id = $1
ORDER BY created_at
