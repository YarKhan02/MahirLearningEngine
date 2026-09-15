SELECT s.id, s.batch_id, s.course_id, s.host_id, s.title, s.class_date, s.status, s.started_at, s.ended_at, c.title, b.batch_name
FROM live_sessions s
JOIN course c ON c.id = s.course_id
JOIN batches b ON b.id = s.batch_id
WHERE s.batch_id = $1 AND s.course_id = $2
ORDER BY s.class_date DESC, s.started_at DESC
