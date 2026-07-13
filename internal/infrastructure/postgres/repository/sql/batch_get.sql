SELECT
    id,
    batch_name,
    start_date,
    end_date,
    capacity,
    days,
    status,
    price
FROM batches
ORDER BY created_at DESC
