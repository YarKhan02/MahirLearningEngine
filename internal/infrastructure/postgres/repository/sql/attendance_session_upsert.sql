INSERT INTO attendance_session (
    id,
    batch_id,
    lesson_date,
    created_by
)
VALUES (
    $1, $2, $3, $4
)
ON CONFLICT ON CONSTRAINT uq_attendance_session
DO UPDATE SET batch_id = EXCLUDED.batch_id
RETURNING id;
