SELECT id, r2_key
FROM attachments
WHERE status = 'pending'
  AND created_at < NOW() - INTERVAL '24 hours'
ORDER BY created_at
LIMIT 500;
