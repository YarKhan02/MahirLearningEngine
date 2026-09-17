DROP TABLE IF EXISTS submission_test_results;
DROP TABLE IF EXISTS assignment_test_cases;

ALTER TABLE assignment_submissions
    DROP COLUMN IF EXISTS auto_score,
    DROP COLUMN IF EXISTS tests_total,
    DROP COLUMN IF EXISTS tests_passed;

ALTER TABLE assignments
    DROP COLUMN IF EXISTS language;
