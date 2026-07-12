SELECT EXISTS (
    SELECT 1
    FROM assignments a
    JOIN lesson l ON l.id = a.lesson_id
    JOIN student_course_access sca ON sca.course_id = l.course_id
    JOIN student_batches sb ON sb.batch_id = sca.batch_id
    WHERE a.id = $2 AND sb.student_id = $1
)
