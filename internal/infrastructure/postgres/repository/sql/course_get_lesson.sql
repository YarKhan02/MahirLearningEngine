SELECT
    id,
    title,
    order_no
FROM lesson
WHERE course_id = $1
ORDER BY order_no ASC;