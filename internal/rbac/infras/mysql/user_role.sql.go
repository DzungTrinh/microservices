// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_role.sql

package mysql

import (
	"context"
	"time"
)

const assignRolesToUser = `-- name: AssignRolesToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE user_id = user_id
`

type AssignRolesToUserParams struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

func (q *Queries) AssignRolesToUser(ctx context.Context, arg AssignRolesToUserParams) error {
	_, err := q.db.ExecContext(ctx, assignRolesToUser, arg.UserID, arg.RoleID)
	return err
}

const listRolesForUser = `-- name: ListRolesForUser :many
SELECT r.id, r.name, r.built_in, r.created_at
FROM roles r
         JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ?
`

type ListRolesForUserRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BuiltIn   bool      `json:"built_in"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) ListRolesForUser(ctx context.Context, userID string) ([]ListRolesForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, listRolesForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRolesForUserRow
	for rows.Next() {
		var i ListRolesForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.BuiltIn,
			&i.CreatedAt,
		); err != nil {
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
