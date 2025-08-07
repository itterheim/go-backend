package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) ListEvents(queryBuilder *EventQueryBuilder) ([]Event, error) {
	where, params := queryBuilder.Build()
	query := fmt.Sprintf(`
		SELECT
			id,
			type,
			timestamp,
			until,
			tags,
			note,
			reference,
			provider_id
		FROM events
		%s
		ORDER BY timestamp ASC
	`, where)
	fmt.Println(query)
	fmt.Println(params)

	rows, err := r.db.Query(context.Background(), query, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.ID, &event.Type, &event.Timestamp, &event.Until, &event.Tags, &event.Note, &event.Reference, &event.ProviderID)
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
			reference,
			provider_id
		FROM events
		WHERE id = $1
	`, id).Scan(
		&event.ID, &event.Type, &event.Timestamp, &event.Until, &event.Tags, &event.Note, &event.Reference, &event.ProviderID,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) CreateEvent(event *Event) (*Event, error) {
	var id int64 = 0
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO events (type, timestamp, until, tags, note, reference, provider_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, event.Type, event.Timestamp, event.Until, event.Tags, event.Note, event.Reference, event.ProviderID).Scan(&id)
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
		    reference = $7,
		    provider_id = $8
		WHERE id = $1
	`, event.ID, event.Type, event.Timestamp, event.Until, event.Tags, event.Note, event.Reference, event.ProviderID)
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
		DELETE FROM events
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

func (r *EventRepository) UsedTags() ([]string, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT DISTINCT unnest(tags) AS unique_tag FROM events;
	`)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for rows.Next() {
		tag := ""
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}

		result = append(result, tag)
	}

	return result, nil
}
