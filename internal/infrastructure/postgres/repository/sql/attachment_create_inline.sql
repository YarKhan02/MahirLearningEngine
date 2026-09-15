INSERT INTO attachments (
    id, r2_key, file_name, content_type, size_bytes,
    resource_type, resource_id, uploaded_by, status,
    verified_content_type, confirmed_at
)
VALUES ($1, $2, $3, $4, $5, 'inline', $6, $7, 'confirmed', $8, NOW())
