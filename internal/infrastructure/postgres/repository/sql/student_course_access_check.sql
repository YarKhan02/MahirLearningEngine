SELECT EXISTS (
    SELECT 1
    FROM student_batches sb
    JOIN student_course_access sca ON sca.batch_id = sb.batch_id
    WHERE sb.student_id = $1 AND sca.course_id = $2
)
