SELECT
    a.id,
    a.lesson_id,
    a.title,
    COALESCE(a.description, ''),
    COALESCE(a.starter_code, ''),
    COALESCE(a.language, ''),
    a.due_date,
    a.total_marks,
    a.created_at,
    s.id,
    s.code,
    s.remarks,
    s.marks,
    s.auto_score,
    s.tests_total,
    s.tests_passed,
    s.status,
    s.submitted_at
FROM assignments a
LEFT JOIN assignment_submissions s
    ON s.assignment_id = a.id AND s.student_id = $2
WHERE a.lesson_id = $1
ORDER BY a.created_at
