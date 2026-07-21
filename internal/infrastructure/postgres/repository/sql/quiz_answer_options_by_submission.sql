SELECT ao.answer_id, ao.option_id
FROM quiz_answer_options ao
JOIN quiz_answers a ON a.id = ao.answer_id
WHERE a.submission_id = $1
