INSERT INTO attendance (
    id,
    session_id,
    student_id,
    status
)
VALUES (
    $1, $2, $3, $4
)
ON CONFLICT ON CONSTRAINT uq_attendance
DO UPDATE SET status = EXCLUDED.status;
