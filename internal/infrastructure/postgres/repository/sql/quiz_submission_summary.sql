SELECT COUNT(*) AS total,
       COUNT(*) FILTER (WHERE status = 'submitted') AS pending,
       COUNT(*) FILTER (WHERE status = 'graded') AS graded
FROM quiz_submissions
WHERE quiz_id = $1
