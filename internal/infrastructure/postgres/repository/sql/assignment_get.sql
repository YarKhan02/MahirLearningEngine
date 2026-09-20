SELECT
    id,
    lesson_id,
    title,
    COALESCE(description, ''),
    COALESCE(starter_code, ''),
    COALESCE(language, ''),
    due_date,
    total_marks,
    created_at
FROM assignments
WHERE id = $1;
