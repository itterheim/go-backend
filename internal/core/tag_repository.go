package core

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db}
}

func (r *TagRepository) ListTags() ([]Tag, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, tag, description, parent_id
		FROM tags
	`)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		err = rows.Scan(
			&tag.ID,
			&tag.Tag,
			&tag.Description,
			&tag.ParentID,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
