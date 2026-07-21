SELECT o.question_id, o.id, o.text, o.is_correct, o.order_no
FROM quiz_options o
JOIN quiz_questions qq ON qq.id = o.question_id
WHERE qq.quiz_id = $1
ORDER BY o.order_no ASC
