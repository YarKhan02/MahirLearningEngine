INSERT INTO submission_test_results
    (id, submission_id, test_case_id, passed, actual_stdout, stderr, timed_out, duration_ms, ordinal)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
