UPDATE batches
SET
    batch_name = $2,
    start_date = $3,
    end_date = $4,
    capacity = $5,
    days = $6,
    status = $7,
    price = $8,
    updated_at = NOW()
WHERE id = $1;
