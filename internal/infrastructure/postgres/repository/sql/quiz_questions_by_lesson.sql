SELECT qq.quiz_id, qq.id, qq.prompt, qq.type, qq.marks, qq.allow_multiple, qq.order_no
FROM quiz_questions qq
JOIN quizzes q ON q.id = qq.quiz_id
WHERE q.lesson_id = $1
ORDER BY qq.order_no ASC
