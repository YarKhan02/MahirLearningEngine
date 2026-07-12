SELECT
    st.id,
    st.full_name,
    st.email,
    a.status
FROM student_batches sb
JOIN students st ON st.id = sb.student_id
LEFT JOIN attendance_session s
    ON s.batch_id = sb.batch_id AND s.lesson_date = $2
LEFT JOIN attendance a
    ON a.session_id = s.id AND a.student_id = st.id
WHERE sb.batch_id = $1
ORDER BY st.full_name
