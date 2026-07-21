SELECT id, lesson_id, title, COALESCE(description, ''), created_at
FROM quizzes
WHERE lesson_id = $1
ORDER BY created_at ASC
