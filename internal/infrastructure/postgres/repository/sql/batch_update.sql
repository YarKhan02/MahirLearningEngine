UPDATE batches
SET
    batch_name = $2,
    start_date = $3,
    end_date = $4,
    capacity = $5,
    status = $6,
    price = $7,
    updated_at = NOW()
WHERE id = $1;
