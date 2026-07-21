SELECT ao.answer_id, ao.option_id
FROM quiz_answer_options ao
JOIN quiz_answers a ON a.id = ao.answer_id
JOIN quiz_submissions s ON s.id = a.submission_id
JOIN quizzes q ON q.id = s.quiz_id
WHERE q.lesson_id = $1 AND s.student_id = $2
