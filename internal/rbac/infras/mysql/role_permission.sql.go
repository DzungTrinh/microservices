// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: role_permission.sql

package mysql

import (
	"context"
	"time"
)

const assignPermissionsToRole = `-- name: AssignPermissionsToRole :exec
INSERT INTO role_permissions (role_id, perm_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE role_id = role_id
`

type AssignPermissionsToRoleParams struct {
	RoleID string `json:"role_id"`
	PermID string `json:"perm_id"`
}

func (q *Queries) AssignPermissionsToRole(ctx context.Context, arg AssignPermissionsToRoleParams) error {
	_, err := q.db.ExecContext(ctx, assignPermissionsToRole, arg.RoleID, arg.PermID)
	return err
}

const listPermissionsForRole = `-- name: ListPermissionsForRole :many
SELECT p.id, p.name, p.created_at
FROM permissions p
         JOIN role_permissions rp ON p.id = rp.perm_id
WHERE rp.role_id = ? AND p.deleted_at IS NULL
`

type ListPermissionsForRoleRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) ListPermissionsForRole(ctx context.Context, roleID string) ([]ListPermissionsForRoleRow, error) {
	rows, err := q.db.QueryContext(ctx, listPermissionsForRole, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPermissionsForRoleRow
	for rows.Next() {
		var i ListPermissionsForRoleRow
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
