INSERT INTO student_course_access (
    id,
    batch_id,
    course_id,
    granted_at,
    granted_by
)
VALUES (
    $1, $2, $3, NOW(), $4
)
ON CONFLICT ON CONSTRAINT uq_batch_course DO NOTHING;
