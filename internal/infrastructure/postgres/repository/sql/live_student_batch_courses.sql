SELECT DISTINCT c.id, c.title
FROM users u
JOIN students s               ON s.username = u.username
JOIN student_batches sb       ON sb.student_id = s.id
JOIN student_course_access sca ON sca.batch_id = sb.batch_id
JOIN course c                 ON c.id = sca.course_id
WHERE u.id = $1 AND sb.batch_id = $2
ORDER BY c.title
