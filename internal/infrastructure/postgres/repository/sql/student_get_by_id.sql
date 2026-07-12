SELECT
    id,
    email,
    full_name,
    phone_number,
    dob,
    gender,
    status
FROM students
WHERE id = $1
