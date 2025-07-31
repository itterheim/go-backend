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
			locations_history.id as l_id, latitude, longitude, accuracy
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
			&location.Extras.ID, &location.Extras.Latitude, &location.Extras.Longitude, &location.Extras.Accuracy,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, location)
	}

	return history, nil
}

func (r *LocationRepository) GetHistory(id int64) (*LocationEvent, error) {
	var data LocationEvent
	err := r.db.QueryRow(context.Background(), `
		SELECT
		    events.id as e_id, type, timestamp, until, tags, note, reference,
			locations_history.id as l_id, latitude, longitude, accuracy
		FROM locations_history
		INNER JOIN events ON locations_history.event_id = events.id
		WHERE locations_history.id = $1
	`, id).Scan(
		&data.ID, &data.Type, &data.Timestamp, &data.Until, &data.Tags, &data.Note, &data.Reference,
		&data.Extras.ID, &data.Extras.Latitude, &data.Extras.Longitude, &data.Extras.Accuracy,
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
		RETURNING id, latitude, longitude, accuracy, event_id
	`, history.Latitude, history.Longitude, history.Accuracy, history.EventID).Scan(
		&result.ID,
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
		SET latitude = $2,
			longitude = $3,
			accuracy = $4
		WHERE id = $1 AND event_id = $5
		RETURNING id, latitude, longitude, accuracy, event_id
	`, history.ID, history.Latitude, history.Longitude, history.Accuracy, history.EventID).Scan(
		&result.ID,
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

func (r *LocationRepository) DeleteHistory(id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM events
		USING locations_history
		WHERE events.id = locations_history.event_id AND locations_history.id = $1
	`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("LocationRepository.DeleteHistory: no rows affected")
	}

	return nil
}
