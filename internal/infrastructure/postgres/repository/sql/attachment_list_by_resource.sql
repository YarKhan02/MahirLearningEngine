SELECT id, r2_key, file_name, content_type, size_bytes,
       resource_type, resource_id, uploaded_by, status, created_at, confirmed_at
FROM attachments
WHERE resource_type = $1
  AND resource_id = $2
  AND status = 'confirmed'
  AND deleted_at IS NULL
ORDER BY created_at DESC
