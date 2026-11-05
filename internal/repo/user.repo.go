package repo

import (
	"context"

	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetUserProfile(ctx context.Context, id int) (model.UserProfile, error) {
	sql := `
	SELECT users.id, users.email, users.full_name, users.created_at, COUNT(links.slug)
	FROM users
	LEFT JOIN links ON links.user_id = users.id
	WHERE users.id = $1
	AND links.deleted_at IS NULL
	GROUP BY users.id`

	var user model.UserProfile
	if err := u.db.QueryRow(ctx, sql, id).Scan(
		&user.Id,
		&user.Email,
		&user.FullName,
		&user.CreatedAt,
		&user.ActiveLink,
	); err != nil {
		return model.UserProfile{}, err
	}

	return user, nil
}
