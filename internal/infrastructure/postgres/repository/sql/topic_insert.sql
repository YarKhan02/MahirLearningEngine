INSERT INTO topics (
    id,
    lesson_id,
    title,
    description,
    content,
    youtube_url,
    order_no
) VALUES (
    $1, $2, $3, $4, $5, $6,
    (SELECT COALESCE(MAX(order_no), 0) + 1 FROM topics WHERE lesson_id = $2)
)
