SELECT
    COUNT(*)                                            AS total,
    COUNT(*) FILTER (WHERE sub.status = 'submitted')    AS submitted,
    COUNT(*) FILTER (WHERE sub.status = 'graded')       AS graded
FROM assignment_submissions sub
WHERE sub.student_id = $1
