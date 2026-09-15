SELECT s.id, s.batch_id, s.course_id, s.host_id, s.title, s.class_date, s.status, s.started_at, s.ended_at, c.title, b.batch_name
FROM live_sessions s
JOIN course c ON c.id = s.course_id
JOIN batches b ON b.id = s.batch_id
WHERE s.id = $1
