SELECT s.id, s.quiz_id, s.student_id, s.status, s.score, s.remarks, s.submitted_at, s.graded_at
FROM quiz_submissions s
JOIN quizzes q ON q.id = s.quiz_id
WHERE q.lesson_id = $1 AND s.student_id = $2
