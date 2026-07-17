SELECT a.id, a.batch_id, a.title, a.description, a.created_at, COALESCE(b.name, '') AS batch_name
FROM announcements a
LEFT JOIN batches b ON a.batch_id = b.id
WHERE a.id = $1;
