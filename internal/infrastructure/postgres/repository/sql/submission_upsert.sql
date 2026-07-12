INSERT INTO assignment_submissions (
    id,
    student_id,
    assignment_id,
    code,
    status,
    submitted_at
)
VALUES (
    $1, $2, $3, $4, 'submitted', NOW()
)
ON CONFLICT ON CONSTRAINT uq_assignment_submission
DO UPDATE SET code = EXCLUDED.code, status = 'submitted', submitted_at = NOW();
