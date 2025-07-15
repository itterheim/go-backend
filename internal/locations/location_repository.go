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

func (r *LocationRepository) ListHistory() ([]GpsHistory, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, latitude, longitude, accuracy, created
		FROM locations_history
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]GpsHistory, 0)
	for rows.Next() {
		gpsHistory := GpsHistory{}

		err := rows.Scan(&gpsHistory.ID, &gpsHistory.Latitude, &gpsHistory.Longitude, &gpsHistory.Accuracy, &gpsHistory.Created)
		if err != nil {
			return nil, err
		}

		history = append(history, gpsHistory)
	}

	return history, nil
}

func (r *LocationRepository) GetHistory(id int64) (*GpsHistory, error) {
	var data GpsHistory
	err := r.db.QueryRow(context.Background(), `
		SELECT id, latitude, longitude, accuracy, created
		FROM locations_history
		WHERE id = $1
	`, id).Scan(&data.ID, &data.Latitude, &data.Longitude, &data.Accuracy, &data.Created)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *LocationRepository) CreateHistory(history *GpsHistory) (*GpsHistory, error) {
	var result GpsHistory
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO locations_history (latitude, longitude, accuracy)
		VALUES ($1, $2, $3)
		RETURNING id, latitude, longitude, accuracy, created
	`, history.Latitude, history.Longitude, history.Accuracy).Scan(
		&result.ID,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
		&result.Created,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *LocationRepository) UpdateHistory(history *GpsHistory) (*GpsHistory, error) {
	var result GpsHistory
	err := r.db.QueryRow(context.Background(), `
		UPDATE locations_history
		SET latitude = $2,
			longitude = $3,
			accuracy = $4
		WHERE id = $1
		RETURNING id, latitude, longitude, accuracy, created
	`, history.ID, history.Latitude, history.Longitude, history.Accuracy).Scan(
		&result.ID,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
		&result.Created,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *LocationRepository) DeleteHistory(id int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM locations_history
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("LocationRepository.DeleteHistory: no rows affected")
	}

	return nil
}
