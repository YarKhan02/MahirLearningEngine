SELECT
    c.id,
    c.title,
    c.level,
    c.duration,
    COALESCE(c.description, ''),
    (SELECT COUNT(*) FROM lesson l WHERE l.course_id = c.id) AS total_lessons,
    (
        SELECT COUNT(*)
        FROM lesson_progress lp
        JOIN lesson l2 ON l2.id = lp.lesson_id
        WHERE l2.course_id = c.id
          AND lp.student_id = s.id
          AND lp.completed
    ) AS completed_lessons
FROM users u
JOIN students s ON s.username = u.username
JOIN student_batches sb ON sb.student_id = s.id
JOIN student_course_access sca ON sca.batch_id = sb.batch_id
JOIN course c ON c.id = sca.course_id
WHERE u.id = $1
ORDER BY c.title
