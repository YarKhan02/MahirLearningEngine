SELECT
    a.id,
    a.batch_id,
    a.title,
    a.description,
    a.created_at,
    b.batch_name
FROM announcements a
JOIN batches b ON b.id = a.batch_id
ORDER BY a.created_at DESC
