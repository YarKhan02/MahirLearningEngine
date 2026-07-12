UPDATE assignment_submissions
SET marks = $2, remarks = $3, status = 'graded'
WHERE id = $1;
