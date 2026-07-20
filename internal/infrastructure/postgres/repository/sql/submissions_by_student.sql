SELECT
    sub.id,
    sub.code,
    sub.remarks,
    sub.marks,
    sub.status,
    sub.submitted_at,
    st.id,
    st.full_name,
    st.email,
    a.id,
    a.title,
    a.total_marks,
    l.id,
    l.title,
    c.id,
    c.title
FROM assignment_submissions sub
JOIN students st ON st.id = sub.student_id
JOIN assignments a ON a.id = sub.assignment_id
JOIN lesson l ON l.id = a.lesson_id
JOIN course c ON c.id = l.course_id
WHERE sub.student_id = $1
  AND ($2 = '' OR sub.status = $2)
ORDER BY sub.submitted_at DESC, sub.id DESC
LIMIT $3 OFFSET $4
