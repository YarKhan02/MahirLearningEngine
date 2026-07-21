SELECT s.id, s.student_id, st.full_name, st.email, s.status, s.submitted_at, s.score
FROM quiz_submissions s
JOIN students st ON st.id = s.student_id
WHERE s.quiz_id = $1
ORDER BY s.submitted_at DESC
