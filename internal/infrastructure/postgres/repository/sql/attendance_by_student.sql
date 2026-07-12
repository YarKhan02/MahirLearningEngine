SELECT
    s.lesson_date,
    a.status,
    b.batch_name
FROM attendance a
JOIN attendance_session s ON s.id = a.session_id
JOIN batches b ON b.id = s.batch_id
WHERE a.student_id = $1
ORDER BY s.lesson_date DESC
