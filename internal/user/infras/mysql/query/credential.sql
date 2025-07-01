-- name: CreateCredential :exec
INSERT INTO credentials (id, user_id, provider, secret_hash, provider_uid, created_at)
VALUES (?, ?, ?, ?, ?, NOW());

-- name: GetCredentialByEmailAndProvider :one
SELECT id, user_id, provider, secret_hash, provider_uid, created_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM credentials
WHERE user_id = (SELECT id FROM users WHERE email = ? AND deleted_at IS NULL) AND provider = ? AND deleted_at IS NULL;