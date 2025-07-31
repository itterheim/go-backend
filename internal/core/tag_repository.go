package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

func (r *TagRepository) GetTag(id int64) (*Tag, error) {
	var tag Tag
	err := r.db.QueryRow(context.Background(), `
		SELECT id, tag, description, parent_id
		FROM tags
		WHERE id = $1
	`, id).Scan(
		&tag.ID,
		&tag.Tag,
		&tag.Description,
		&tag.ParentID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *TagRepository) CreateTag(data *Tag) (*Tag, error) {
	var id int64 = 0
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO tags (tag, description, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, data.Tag, data.Description, data.ParentID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.GetTag(id)
}

func (r *TagRepository) UpdateTag(data *Tag) (*Tag, error) {
	var id int64
	err := r.db.QueryRow(context.Background(), `
		UPDATE tags
		SET tag = $2,
			description = $3,
			parent_id = $4
		WHERE id = $1
		RETURNING id
	`, &data.ID, &data.Tag, &data.Description, &data.ParentID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.GetTag(id)
}

func (r *TagRepository) DeleteTag(id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM tags
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("TagRepository.DeleteTag: no rows affected")
	}

	return nil
}
