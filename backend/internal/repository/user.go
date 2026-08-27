package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
)

var ErrUserNotFound = errors.New("repository: ไม่พบผู้ใช้")

type User struct {
	ID       uuid.UUID
	Username string
	Role     string
	OrgID    uuid.UUID
	OrgName  string
}

func (r *Repo) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	const q = `
SELECT u.id, u.username, u.role, u.organization_id, o.name_th
FROM users u
JOIN organizations o ON o.id = u.organization_id
WHERE u.id = $1 AND u.is_active AND o.is_active`

	var u User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Role, &u.OrgID, &u.OrgName)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (r *Repo) ListUsers(ctx context.Context) ([]dto.DevUser, error) {
	const q = `
SELECT u.id, u.username, u.role, o.name_th
FROM users u
JOIN organizations o ON o.id = u.organization_id
WHERE u.is_active
ORDER BY u.username`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []dto.DevUser{}
	for rows.Next() {
		var u dto.DevUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.OrgName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
