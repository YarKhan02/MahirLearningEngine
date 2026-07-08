INSERT INTO course (
    id,
    title,
    level,
    duration,
    description
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING
    id,
    title,
    level,
    duration,
    description,
    is_active;