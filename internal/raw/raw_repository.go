package raw

import (
	"backend/internal/core"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const RawTable string = "raw"

type RawRepository struct {
	db *pgxpool.Pool
}

func NewRawRepository(db *pgxpool.Pool) *RawRepository {
	return &RawRepository{db}
}

func (r *RawRepository) ListRawEvents(queryBuilder *core.EventQueryBuilder) ([]RawEvent, error) {
	where, params := queryBuilder.Build()
	query := fmt.Sprintf(`
		SELECT
		    events.id as e_id, type, timestamp, until, tags, note, reference,
			event_id, data
		FROM raw
		INNER JOIN events ON raw.event_id = events.id
		%s
		ORDER BY timestamp ASC
	`, where)

	rows, err := r.db.Query(context.Background(), query, params...)
	if err != nil {
		return nil, fmt.Errorf("RawRepository.ListRawEvents: %v", err)
	}

	result := make([]RawEvent, 0)
	for rows.Next() {
		var data RawEvent

		err := rows.Scan(
			&data.ID, &data.Type, &data.Timestamp, &data.Until, &data.Tags, &data.Note, &data.Reference,
			&data.Extras.EventID, &data.Extras.Data,
		)
		if err != nil {
			return nil, fmt.Errorf("RawRepository.ListRawEvents - failed to parse row: %v", err)
		}

		result = append(result, data)
	}

	return result, nil
}

func (r *RawRepository) GetRawEvent(eventID int64) (*RawEvent, error) {
	result := RawEvent{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
  			events.id as e_id, type, timestamp, until, tags, note, reference,
     		event_id, data
		FROM raw
		INNER JOIN events ON raw.event_id = events.id
		WHERE
		WHERE events.id = $1
	`, eventID).Scan(
		&result.ID, &result.Type, &result.Timestamp, &result.Until, &result.Tags, &result.Note, &result.Reference,
		&result.Extras.EventID, &result.Extras.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("RawRepository.GetRawEvent: %v", err)
	}

	return &result, nil
}

func (r *RawRepository) CreateRaw(data *Raw) (*Raw, error) {
	var result Raw
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO raw (event_id, data)
		VALUES ($1, $2)
		RETURNING event_id, data
	`, data.EventID, data.Data).Scan(&result.EventID, &result.Data)
	if err != nil {
		return nil, fmt.Errorf("RawRepository.CreateRaw: %v", err)
	}

	return &result, nil
}

func (r *RawRepository) UpdateRaw(data *Raw) (*Raw, error) {
	var result Raw
	err := r.db.QueryRow(context.Background(), `
		UPDATE raw
		SET data = $2
		WHERE event_id = $1
		RETURNING event_id, data
	`, data.EventID, data.Data).Scan(&result.EventID, &result.Data)
	if err != nil {
		return nil, fmt.Errorf("RawRepository.UpdateRaw: %v", err)
	}

	return &result, nil
}

func (r *RawRepository) DeleteRawEvent(event_id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM events
		USING raw
		WHERE events.id = raw.event_id AND raw.event_id = $1
	`, event_id)
	if err != nil {
		return fmt.Errorf("RawRepository.DelteRaw: %v", err)
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("RawRepository.DelteRaw: no rows affected")
	}

	return nil
}
