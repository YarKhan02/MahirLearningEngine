SELECT id, question_id, COALESCE(answer_text, ''), awarded_marks
FROM quiz_answers
WHERE submission_id = $1
