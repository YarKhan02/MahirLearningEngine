SELECT test_case_id, passed, actual_stdout, stderr, timed_out, duration_ms, ordinal
FROM submission_test_results
WHERE submission_id = $1
ORDER BY ordinal;
