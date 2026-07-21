SELECT id, prompt, type, marks, allow_multiple, order_no
FROM quiz_questions
WHERE quiz_id = $1
ORDER BY order_no ASC
