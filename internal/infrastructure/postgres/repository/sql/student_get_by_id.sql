SELECT
    id,
    email,
    username,
    full_name,
    phone_number,
    dob,
    gender,
    status
FROM students
WHERE id = $1
