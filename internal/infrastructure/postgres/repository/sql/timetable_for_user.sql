SELECT
    t.id,
    t.batch_id,
    t.course_id,
    t.weekdays,
    t.start_time,
    t.end_time,
    c.title,
    b.batch_name,
    b.start_date,
    b.end_date
FROM users u
JOIN students s ON s.username = u.username
JOIN student_batches sb ON sb.student_id = s.id
JOIN timetable t ON t.batch_id = sb.batch_id
JOIN course c ON c.id = t.course_id
JOIN batches b ON b.id = t.batch_id
WHERE u.id = $1
