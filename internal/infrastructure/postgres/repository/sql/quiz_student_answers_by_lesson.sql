SELECT a.id, a.submission_id, a.question_id, COALESCE(a.answer_text, ''), a.awarded_marks
FROM quiz_answers a
JOIN quiz_submissions s ON s.id = a.submission_id
JOIN quizzes q ON q.id = s.quiz_id
WHERE q.lesson_id = $1 AND s.student_id = $2
