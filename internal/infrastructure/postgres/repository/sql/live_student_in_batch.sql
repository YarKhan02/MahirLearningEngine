SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN students s         ON s.username = u.username
    JOIN student_batches sb ON sb.student_id = s.id
    WHERE u.id = $1 AND sb.batch_id = $2
)
