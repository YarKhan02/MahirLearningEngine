INSERT INTO student_batches (
    id,
    student_id,
    batch_id,
    enrolled_at
)
VALUES (
    $1, $2, $3, NOW()
);
