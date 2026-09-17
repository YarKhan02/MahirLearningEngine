UPDATE assignment_submissions
SET auto_score = $2, tests_total = $3, tests_passed = $4
WHERE id = $1;
