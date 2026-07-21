SELECT id, r2_key, file_name, content_type, size_bytes,
       resource_type, resource_id, uploaded_by, status, created_at, confirmed_at
FROM attachments
WHERE id = $1
  AND deleted_at IS NULL
