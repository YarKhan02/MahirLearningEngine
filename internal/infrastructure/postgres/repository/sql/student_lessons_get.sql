SELECT
    l.id,
    l.title,
    l.order_no,
    COALESCE(lp.completed, FALSE) AS completed,
    lp.completed_at
FROM lesson l
LEFT JOIN lesson_progress lp
    ON lp.lesson_id = l.id AND lp.student_id = $2
WHERE l.course_id = $1
ORDER BY l.order_no
