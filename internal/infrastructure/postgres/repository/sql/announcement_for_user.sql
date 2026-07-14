SELECT
    a.id,
    a.batch_id,
    a.title,
    a.description,
    a.created_at,
    b.batch_name
FROM users u
JOIN students s ON s.username = u.username
JOIN student_batches sb ON sb.student_id = s.id
JOIN announcements a ON a.batch_id = sb.batch_id
JOIN batches b ON b.id = a.batch_id
WHERE u.id = $1
ORDER BY a.created_at DESC
