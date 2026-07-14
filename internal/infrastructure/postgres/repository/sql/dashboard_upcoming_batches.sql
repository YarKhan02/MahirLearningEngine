SELECT
    b.id,
    b.batch_name,
    b.start_date,
    b.price,
    b.capacity,
    (SELECT COUNT(*) FROM student_batches sb WHERE sb.batch_id = b.id) AS enrolled
FROM batches b
WHERE b.start_date >= CURRENT_DATE
  AND b.start_date < CURRENT_DATE + INTERVAL '7 days'
ORDER BY b.start_date
