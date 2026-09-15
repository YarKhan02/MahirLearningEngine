SELECT id, session_id, storage_key, image_url, version, created_by, created_at
FROM whiteboard_snapshots
WHERE session_id = $1
ORDER BY version DESC
