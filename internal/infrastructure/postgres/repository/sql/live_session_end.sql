UPDATE live_sessions
SET status = 'ended', ended_at = NOW()
WHERE id = $1 AND host_id = $2 AND status = 'live'
