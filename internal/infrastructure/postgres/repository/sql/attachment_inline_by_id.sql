SELECT r2_key, content_type
FROM attachments
WHERE id = $1
  AND resource_type = 'inline'
  AND deleted_at IS NULL
