SELECT EXISTS (
    SELECT 1
    FROM live_sessions ls
    JOIN student_batches sb ON sb.batch_id = ls.batch_id
    JOIN students s         ON s.id = sb.student_id
    JOIN users u            ON u.username = s.username
    JOIN student_course_access sca
        ON sca.batch_id = ls.batch_id AND sca.course_id = ls.course_id
    WHERE u.id = $1 AND ls.id = $2
)
