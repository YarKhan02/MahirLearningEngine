UPDATE quiz_submissions
SET status = 'graded',
    remarks = $1,
    graded_at = NOW(),
    score = (SELECT COALESCE(SUM(awarded_marks), 0) FROM quiz_answers WHERE submission_id = $2)
WHERE id = $2
