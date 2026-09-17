SELECT id, stdin, expected_stdout, weight, ordinal
FROM assignment_test_cases
WHERE assignment_id = $1
ORDER BY ordinal;
