SELECT DISTINCT b.id, b.batch_name
FROM users u
JOIN students s          ON s.username = u.username
JOIN student_batches sb  ON sb.student_id = s.id
JOIN batches b           ON b.id = sb.batch_id
WHERE u.id = $1
ORDER BY b.batch_name
