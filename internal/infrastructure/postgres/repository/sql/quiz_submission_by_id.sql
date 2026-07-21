SELECT id, quiz_id, student_id, status, score, remarks, submitted_at, graded_at
FROM quiz_submissions
WHERE id = $1
