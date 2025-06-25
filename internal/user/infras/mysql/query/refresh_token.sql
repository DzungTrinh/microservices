-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token, user_agent, ip_address, expires_at, revoked)
VALUES (?, ?, ?, ?, ?, ?, 0);

-- name: GetRefreshToken :one
SELECT id, user_id, token, user_agent, ip_address, created_at, expires_at, revoked
FROM refresh_tokens
WHERE token = ? AND expires_at > NOW() AND revoked = 0;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = 1
WHERE token = ?;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW() OR revoked = 1;