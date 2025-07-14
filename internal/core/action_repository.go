package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActionRepository struct {
	db *pgxpool.Pool
}

func NewActionRepository(db *pgxpool.Pool) *ActionRepository {

	return &ActionRepository{db: db}
}

func (r *ActionRepository) ListActions() ([]ActionResponse, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, event_id, reference, reference_id, tags, note
		FROM actions
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actions := make([]ActionResponse, 0)
	for rows.Next() {
		action := ActionResponse{}
		reference := ActionReference{}

		err = rows.Scan(&action.ID, &action.EventID, &reference.Table, &reference.ID, &action.Tags, &action.Note)
		if err != nil {
			return nil, err
		}

		if reference.Table != nil {
			action.Reference = &reference
		}

		actions = append(actions, action)
	}

	return actions, nil

}

func (r *ActionRepository) GetAction(id int64) (*ActionResponse, error) {
	var action ActionResponse
	var reference ActionReference
	err := r.db.QueryRow(context.Background(), `
		SELECT id, event_id, reference, reference_id, tags, note
		FROM actions
		WHERE id = $1
	`, id).Scan(&action.ID, &action.EventID, &reference.Table, &reference.ID, &action.Tags, &action.Note)
	if err != nil {
		return nil, err
	}

	if reference.Table != nil {
		action.Reference = &reference
	}

	return &action, nil
}

func (r *ActionRepository) CreateAction(action *CreateActionRequest) (*ActionResponse, error) {
	var reference *string
	var referenceId *int64
	if action.Reference != nil {
		reference = action.Reference.Table
		referenceId = action.Reference.ID
	}

	var id int64
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO actions (event_id, reference, reference_id, tags, note)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, action.EventID, reference, referenceId, action.Tags, action.Note).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.GetAction(id)
}

func (r *ActionRepository) UpdateAction(action *UpdateActionRequest) (*ActionResponse, error) {
	var reference *string
	var referenceId *int64
	if action.Reference != nil {
		reference = action.Reference.Table
		referenceId = action.Reference.ID
	}

	cmd, err := r.db.Exec(context.Background(), `
		UPDATE actions
		SET event_id = $2,
			reference = $3,
			reference_id = $4,
			tags = $5,
			note = $6
		WHERE id = $1
	`, action.ID, action.EventID, reference, referenceId, action.Tags, action.Note)
	if err != nil {
		return nil, err
	}

	if cmd.RowsAffected() == 0 {
		return nil, errors.New("ActionRepository.UpdateAction: no rows updated")
	}

	if cmd.RowsAffected() > 1 {
		return nil, errors.New("ActionRepository.UpdateAction: too many rows updated")
	}

	return r.GetAction(action.ID)
}

func (r *ActionRepository) DeleteAction(id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM actions
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("Delete: no rows affected")
	}

	return nil
}
