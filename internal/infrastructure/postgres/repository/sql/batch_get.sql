SELECT
    id,
    batch_name,
    start_date,
    end_date,
    capacity,
    status,
    price
FROM batches
ORDER BY created_at DESC
