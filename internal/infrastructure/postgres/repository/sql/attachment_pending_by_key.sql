SELECT id, r2_key, file_name, content_type, size_bytes,
       resource_type, resource_id, uploaded_by, status, created_at, confirmed_at
FROM attachments
WHERE r2_key = $1
  AND uploaded_by = $2
  AND status = 'pending'
  AND deleted_at IS NULL
