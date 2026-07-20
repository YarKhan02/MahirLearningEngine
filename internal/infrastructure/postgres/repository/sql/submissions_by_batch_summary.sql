SELECT
    COUNT(*)                                            AS total,
    COUNT(*) FILTER (WHERE sub.status = 'submitted')    AS submitted,
    COUNT(*) FILTER (WHERE sub.status = 'graded')       AS graded
FROM assignment_submissions sub
JOIN students st ON st.id = sub.student_id
JOIN student_batches sb ON sb.student_id = st.id AND sb.batch_id = $1
JOIN assignments a ON a.id = sub.assignment_id
JOIN lesson l ON l.id = a.lesson_id
JOIN course c ON c.id = l.course_id
WHERE ($2 = '' OR st.full_name ILIKE '%' || $2 || '%'
                OR st.email ILIKE '%' || $2 || '%'
                OR a.title ILIKE '%' || $2 || '%'
                OR l.title ILIKE '%' || $2 || '%'
                OR c.title ILIKE '%' || $2 || '%')
