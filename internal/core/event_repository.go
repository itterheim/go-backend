package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) ListEvents() ([]Event, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
			id,
			type,
			timestamp,
			until,
			status,
			tags,
			note
		FROM events
		ORDER BY timestamp ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.ID, &event.Type, &event.Timestamp, &event.Until, &event.Status, &event.Tags, &event.Note)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetEvent(id int64) (*Event, error) {
	event := Event{}

	err := r.db.QueryRow(context.Background(), `
		SELECT
		    id,
		    type,
		    timestamp,
		    until,
		    tags,
		    note,
		    status
		FROM events
		WHERE id = $1
	`, id).Scan(
		&event.ID, &event.Type, &event.Timestamp, &event.Until, &event.Tags, &event.Note, &event.Status,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) CreateEvent(event *Event) (*Event, error) {
	var id int64 = 0
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO events (type, timestamp, until, tags, note, status, provider_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, event.Type, event.Timestamp, event.Until, event.Tags, event.Note, event.Status, event.ProviderID, event.UserID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.GetEvent(id)
}

func (r *EventRepository) UpdateEvent(event *Event) (*Event, error) {
	cmd, err := r.db.Exec(context.Background(), `
		UPDATE events
		SET type = $2,
		    timestamp = $3,
		    until = $4,
		    tags = $5,
		    note = $6,
		    status = $7,
		    provider_id = $8,
		    user_id = $9
		WHERE id = $1
	`, event.ID, event.Type, event.Timestamp, event.Until, event.Tags, event.Note, event.Status, event.ProviderID, event.UserID)
	if err != nil {
		return nil, err
	}

	if cmd.RowsAffected() == 0 {
		return nil, errors.New("Update: no rows updated")
	}

	if cmd.RowsAffected() > 1 {
		return nil, errors.New("Update: too many rows updated")
	}

	return r.GetEvent(event.ID)
}

func (r *EventRepository) DeleteEvent(id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FORM events
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
