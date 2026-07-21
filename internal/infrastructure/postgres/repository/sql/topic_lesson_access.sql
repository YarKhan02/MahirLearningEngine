SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN students s ON s.username = u.username
    JOIN student_batches sb ON sb.student_id = s.id
    JOIN student_course_access sca ON sca.batch_id = sb.batch_id
    JOIN lesson l ON l.course_id = sca.course_id
    WHERE u.id = $1 AND l.id = $2
)
