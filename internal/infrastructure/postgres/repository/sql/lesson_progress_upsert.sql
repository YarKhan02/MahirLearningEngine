INSERT INTO lesson_progress (
    id,
    student_id,
    lesson_id,
    completed,
    completed_at
)
VALUES (
    $1, $2, $3, $4, NOW()
)
ON CONFLICT ON CONSTRAINT uq_lesson_progress
DO UPDATE SET completed = EXCLUDED.completed, completed_at = NOW();
