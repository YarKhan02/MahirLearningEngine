SELECT s.id, s.batch_id, s.course_id, s.host_id, s.title, s.class_date, s.status, s.started_at, s.ended_at, c.title, b.batch_name
FROM live_sessions s
JOIN course c ON c.id = s.course_id
JOIN batches b ON b.id = s.batch_id
JOIN student_batches sb ON sb.batch_id = s.batch_id
JOIN students st       ON st.id = sb.student_id
JOIN users u           ON u.username = st.username
WHERE u.id = $1 AND s.status = 'live'
LIMIT 1
