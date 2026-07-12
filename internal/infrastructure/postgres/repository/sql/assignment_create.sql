INSERT INTO assignments (
    id,
    lesson_id,
    title,
    description,
    starter_code,
    due_date,
    total_marks,
    created_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
);
