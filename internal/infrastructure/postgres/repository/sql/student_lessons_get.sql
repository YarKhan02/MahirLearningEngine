SELECT
    l.id,
    l.title,
    COALESCE(l.description, ''),
    l.order_no,
    COALESCE(l.youtube_url, ''),
    COALESCE(l.content, ''),
    COALESCE(lp.completed, FALSE) AS completed,
    lp.completed_at
FROM lesson l
LEFT JOIN lesson_progress lp
    ON lp.lesson_id = l.id AND lp.student_id = $2
WHERE l.course_id = $1
ORDER BY l.order_no
