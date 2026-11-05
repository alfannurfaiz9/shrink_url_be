package repo

import (
	"context"

	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinksRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (l *LinkRepo) CheckExistingLinks(ctx context.Context, slug string) (bool, error) {
	sql := `
	SELECT EXISTS (
		SELECT 1 from links 
		WHERE slug = $1
	)
	`
	var exist bool
	if err := l.db.QueryRow(ctx, sql, slug).Scan(
		&exist,
	); err != nil {
		return false, err
	}

	return exist, nil
}

func (l *LinkRepo) CreateLink(ctx context.Context, body model.Link) error {
	sql := `
	INSERT INTO links(user_id, original_url, slug)
	VALUES($1, $2, $3)
	`
	args := []any{body.UserId, body.OriginalUrl, body.Slug}

	_, err := l.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepo) GetUserLink(ctx context.Context, userId int) ([]model.Link, error) {
	sql := `
	SELECT id, slug, original_url, created_at
	FROM links
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC`
	rows, err := l.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	links := make([]model.Link, 0)
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(
			&link.Id,
			&link.Slug,
			&link.OriginalUrl,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
