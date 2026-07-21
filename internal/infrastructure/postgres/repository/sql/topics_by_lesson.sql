SELECT
    id,
    lesson_id,
    title,
    COALESCE(description, ''),
    COALESCE(content, ''),
    COALESCE(youtube_url, ''),
    order_no
FROM topics
WHERE lesson_id = $1
ORDER BY order_no ASC
