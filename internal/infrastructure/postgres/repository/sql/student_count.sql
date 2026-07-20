SELECT COUNT(*)
FROM students s
LEFT JOIN student_batches sb ON sb.student_id = s.id
LEFT JOIN batches b ON b.id = sb.batch_id
WHERE ($1 = '' OR s.full_name ILIKE '%' || $1 || '%' OR s.email ILIKE '%' || $1 || '%' OR s.username ILIKE '%' || $1 || '%')
