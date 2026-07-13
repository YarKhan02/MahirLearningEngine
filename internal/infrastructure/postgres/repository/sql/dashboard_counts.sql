SELECT
    (SELECT COUNT(*) FROM students)                                        AS total_students,
    (SELECT COUNT(*) FROM students WHERE status = 'active')                AS active_students,
    (SELECT COUNT(*) FROM students WHERE status = 'pending')               AS pending_students,
    (SELECT COUNT(*) FROM assignment_submissions WHERE status = 'submitted') AS pending_submissions
