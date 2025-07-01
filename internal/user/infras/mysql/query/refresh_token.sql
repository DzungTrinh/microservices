-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token, user_agent, ip_address, created_at, expires_at, revoked)
VALUES (?, ?, ?, ?, ?, NOW(), ?, ?);

-- name: GetRefreshToken :one
SELECT id, user_id, token, user_agent, ip_address, created_at, expires_at, revoked, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM refresh_tokens
WHERE token = ? AND deleted_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = 1, deleted_at = NOW()
WHERE token = ? AND deleted_at IS NULL;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW() OR revoked = 1;
