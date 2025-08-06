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

func (r *TagRepository) ListTags(private bool) ([]Tag, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT tag, description, parent, private
		FROM tags
		WHERE $1 OR private
	`, private)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		err = rows.Scan(
			&tag.Tag,
			&tag.Description,
			&tag.Parent,
			&tag.Private,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) GetTag(tag string) (*Tag, error) {
	var result Tag
	err := r.db.QueryRow(context.Background(), `
		SELECT tag, description, parent, private
		FROM tags
		WHERE tag = $1
	`, tag).Scan(
		&result.Tag,
		&result.Description,
		&result.Parent,
		&result.Private,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *TagRepository) CreateTag(data *Tag) (*Tag, error) {
	var tag string
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO tags (tag, description, parent, private)
		VALUES ($1, $2, $3, $4)
		RETURNING tag
	`, data.Tag, data.Description, data.Parent, data.Private).Scan(&tag)
	if err != nil {
		return nil, err
	}

	return r.GetTag(tag)
}

func (r *TagRepository) UpdateTag(data *Tag, rename *string) (*Tag, error) {
	newTag := data.Tag
	if rename != nil {
		newTag = *rename
	}
	var tag string
	err := r.db.QueryRow(context.Background(), `
		UPDATE tags
		SET tag = $2,
			description = $3,
			parent = $4
		WHERE tag = $1
		RETURNING tag
	`, &data.Tag, newTag, &data.Description, &data.Parent).Scan(&tag)
	if err != nil {
		return nil, err
	}

	return r.GetTag(tag)
}

func (r *TagRepository) DeleteTag(tag string) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM tags
		WHERE tag = $1
	`, tag)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("TagRepository.DeleteTag: no rows affected")
	}

	return nil
}
