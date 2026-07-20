SELECT COUNT(*)
FROM assignment_submissions sub
WHERE sub.student_id = $1
  AND ($2 = '' OR sub.status = $2)
