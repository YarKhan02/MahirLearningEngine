SELECT s.quiz_id,
       COUNT(*) AS total,
       COUNT(*) FILTER (WHERE s.status = 'submitted') AS pending,
       COUNT(*) FILTER (WHERE s.status = 'graded') AS graded
FROM quiz_submissions s
JOIN quizzes q ON q.id = s.quiz_id
WHERE q.lesson_id = $1
GROUP BY s.quiz_id
