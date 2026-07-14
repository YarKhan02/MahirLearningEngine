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
FROM timetable t
JOIN course c ON c.id = t.course_id
JOIN batches b ON b.id = t.batch_id
WHERE t.batch_id = $1
ORDER BY t.start_time
