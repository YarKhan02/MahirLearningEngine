INSERT INTO attachments (
    id, r2_key, file_name, content_type, size_bytes,
    resource_type, resource_id, uploaded_by, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
