SELECT
    id,
    title,
    description,
    order_no,
    youtube_url,
    content
FROM lesson
WHERE course_id = $1
ORDER BY order_no ASC;