INSERT INTO lesson (
    id,
    course_id,
    title,
    description,
    order_no,
    youtube_url,
    content
) VALUES ($1, $2, $3, $4, $5, $6, $7)