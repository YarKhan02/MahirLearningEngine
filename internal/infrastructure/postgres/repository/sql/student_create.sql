INSERT INTO students (
    id,
    email,
    full_name,
    phone_number,
    dob,
    gender,
    status,
    created_at,
    updated_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
);
