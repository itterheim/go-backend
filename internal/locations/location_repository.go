package locations

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const LocationGPSHistoryTable string = "locations_history"
const LocationPlacesTable string = "locations_places"

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db}
}

func (r *LocationRepository) ListHistory() ([]LocationEvent, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
		    events.id as e_id, type, timestamp, until, tags, note, reference,
			event_id, latitude, longitude, accuracy
		FROM locations_history
		INNER JOIN events ON locations_history.event_id = events.id
		ORDER BY timestamp ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]LocationEvent, 0)
	for rows.Next() {
		location := LocationEvent{}

		err := rows.Scan(
			&location.ID, &location.Type, &location.Timestamp, &location.Until, &location.Tags, &location.Note, &location.Reference,
			&location.Extras.EventID, &location.Extras.Latitude, &location.Extras.Longitude, &location.Extras.Accuracy,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, location)
	}

	return history, nil
}

func (r *LocationRepository) GetHistory(eventId int64) (*LocationEvent, error) {
	var data LocationEvent
	err := r.db.QueryRow(context.Background(), `
		SELECT
		    events.id as e_id, type, timestamp, until, tags, note, reference,
			event_id, latitude, longitude, accuracy
		FROM locations_history
		INNER JOIN events ON locations_history.event_id = events.id
		WHERE events.id = $1
	`, eventId).Scan(
		&data.ID, &data.Type, &data.Timestamp, &data.Until, &data.Tags, &data.Note, &data.Reference,
		&data.Extras.EventID, &data.Extras.Latitude, &data.Extras.Longitude, &data.Extras.Accuracy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *LocationRepository) CreateHistory(history *Location) (*Location, error) {
	var result Location
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO locations_history (latitude, longitude, accuracy, event_id)
		VALUES ($1, $2, $3, $4)
		RETURNING latitude, longitude, accuracy, event_id
	`, history.Latitude, history.Longitude, history.Accuracy, history.EventID).Scan(
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
		&result.EventID,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *LocationRepository) UpdateHistory(history *Location) (*Location, error) {
	var result Location
	err := r.db.QueryRow(context.Background(), `
		UPDATE locations_history
		SET latitude = $1,
			longitude = $2,
			accuracy = $3
		WHERE event_id = $4
		RETURNING id, latitude, longitude, accuracy, event_id
	`, history.Latitude, history.Longitude, history.Accuracy, history.EventID).Scan(
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
		&result.EventID,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *LocationRepository) DeleteHistory(event_id int64) error {
	// NOTE: this function is actually not necessary because the event can be deleted directly and history will be deleted thanks to the db constraint
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM events
		USING locations_history
		WHERE events.id = locations_history.event_id AND locations_history.event_id = $1
	`, event_id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("LocationRepository.DeleteHistory: no rows affected")
	}

	return nil
}
