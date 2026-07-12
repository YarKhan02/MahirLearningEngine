INSERT INTO batches (
    id,
    batch_name,
    start_date,
    end_date,
    capacity,
    days,
    status,
    created_at,
    updated_at
) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
);