UPDATE attachments
SET status = 'confirmed',
    confirmed_at = NOW(),
    size_bytes = $3
WHERE r2_key = $1
  AND uploaded_by = $2
  AND deleted_at IS NULL
RETURNING id, r2_key, file_name, content_type, size_bytes,
          resource_type, resource_id, uploaded_by, status, created_at, confirmed_at
