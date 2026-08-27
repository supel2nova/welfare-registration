package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrOrgNotFound = errors.New("repository: ไม่พบหน่วยรับลงทะเบียน")

func (r *Repo) OrgName(ctx context.Context, id uuid.UUID) (string, error) {
	var name string
	err := r.pool.QueryRow(ctx, `SELECT name_th FROM organizations WHERE id = $1 AND is_active`, id).Scan(&name)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrOrgNotFound
	}
	return name, err
}
