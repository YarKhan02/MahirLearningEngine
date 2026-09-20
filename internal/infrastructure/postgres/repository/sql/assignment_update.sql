UPDATE assignments
SET title = $2,
    description = $3,
    starter_code = $4,
    language = $5,
    due_date = $6,
    total_marks = $7
WHERE id = $1;
