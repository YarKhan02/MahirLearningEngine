SELECT id, lesson_id, title, COALESCE(description, ''), created_at
FROM quizzes
WHERE id = $1
