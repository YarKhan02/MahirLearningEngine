DELETE FROM refresh_tokens
WHERE expires_at < NOW()
   OR (revoked AND revoked_at IS NOT NULL AND revoked_at < NOW() - INTERVAL '7 days');
