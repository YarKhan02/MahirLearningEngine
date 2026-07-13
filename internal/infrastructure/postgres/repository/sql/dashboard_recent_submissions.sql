SELECT
    sub.id,
    st.full_name,
    a.title,
    c.title,
    sub.status,
    sub.submitted_at
FROM assignment_submissions sub
JOIN students st ON st.id = sub.student_id
JOIN assignments a ON a.id = sub.assignment_id
JOIN lesson l ON l.id = a.lesson_id
JOIN course c ON c.id = l.course_id
ORDER BY sub.submitted_at DESC
LIMIT 5
