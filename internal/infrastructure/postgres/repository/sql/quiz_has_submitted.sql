SELECT EXISTS(SELECT 1 FROM quiz_submissions WHERE quiz_id = $1 AND student_id = $2)
