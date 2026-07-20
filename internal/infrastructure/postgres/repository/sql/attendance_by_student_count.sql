SELECT COUNT(*)
FROM attendance a
JOIN attendance_session s ON s.id = a.session_id
WHERE a.student_id = $1
