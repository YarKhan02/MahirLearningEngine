SELECT
    sca.id,
    sca.course_id,
    c.title,
    c.level,
    sca.granted_at
FROM student_course_access sca
JOIN course c ON c.id = sca.course_id
WHERE sca.batch_id = $1
ORDER BY sca.granted_at DESC
