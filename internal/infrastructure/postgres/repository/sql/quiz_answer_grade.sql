UPDATE quiz_answers SET awarded_marks = $1
WHERE submission_id = $2 AND question_id = $3
