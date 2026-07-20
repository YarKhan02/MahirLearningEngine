SELECT
    COUNT(*) FILTER (WHERE a.status = 'present')                     AS present,
    COUNT(*) FILTER (WHERE a.status = 'absent')                      AS absent,
    COUNT(*)                                                         AS total,
    MAX(a.status) FILTER (WHERE s.lesson_date = CURRENT_DATE)        AS today_status
FROM attendance a
JOIN attendance_session s ON s.id = a.session_id
WHERE a.student_id = $1
