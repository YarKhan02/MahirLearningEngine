SELECT EXISTS (
    SELECT 1
    FROM lesson l
    JOIN student_course_access sca ON sca.course_id = l.course_id
    JOIN student_batches sb ON sb.batch_id = sca.batch_id
    WHERE l.id = $2 AND sb.student_id = $1
)
