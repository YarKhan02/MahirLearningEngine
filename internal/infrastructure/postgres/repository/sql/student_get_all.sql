SELECT
    s.id,
    s.email,
    s.full_name,
    s.phone_number,
    s.dob,
    s.gender,
    s.status,
    sb.batch_id,
    b.batch_name,
    EXISTS (SELECT 1 FROM users u WHERE u.email = s.email) AS has_account
FROM students s
LEFT JOIN student_batches sb ON sb.student_id = s.id
LEFT JOIN batches b ON b.id = sb.batch_id
WHERE ($1 = '' OR s.full_name ILIKE '%' || $1 || '%' OR s.email ILIKE '%' || $1 || '%')
ORDER BY s.created_at DESC
