SELECT COALESCE(MAX(version), 0) + 1 FROM whiteboard_snapshots WHERE session_id = $1
