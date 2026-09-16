UPDATE live_sessions
SET status = 'ended', ended_at = NOW()
WHERE status = 'live'
