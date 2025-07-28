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

func (r *LocationRepository) ListHistory(userID int64) ([]GpsHistory, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, timestamp, latitude, longitude, accuracy
		FROM locations_history
		WHERE user_id = $1
		ORDER BY timestamp ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]GpsHistory, 0)
	for rows.Next() {
		gpsHistory := GpsHistory{}

		err := rows.Scan(&gpsHistory.ID, &gpsHistory.Timestamp, &gpsHistory.Latitude, &gpsHistory.Longitude, &gpsHistory.Accuracy)
		if err != nil {
			return nil, err
		}

		history = append(history, gpsHistory)
	}

	return history, nil
}

func (r *LocationRepository) GetHistory(id, userID int64) (*GpsHistory, error) {
	var data GpsHistory
	err := r.db.QueryRow(context.Background(), `
		SELECT id, timestamp, latitude, longitude, accuracy
		FROM locations_history
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&data.ID, &data.Timestamp, &data.Latitude, &data.Longitude, &data.Accuracy)

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
		INSERT INTO locations_history (timestamp, latitude, longitude, accuracy, provider_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, timestamp, latitude, longitude, accuracy
	`, history.Timestamp, history.Latitude, history.Longitude, history.Accuracy, history.ProviderID, history.UserID).Scan(
		&result.ID,
		&result.Timestamp,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
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
		SET timestamp = $2,
		    latitude = $3,
			longitude = $4,
			accuracy = $5,
			provider_id = $6
		WHERE id = $1 AND user_id = $7
		RETURNING id, timestamp, latitude, longitude, accuracy
	`, history.ID, history.Timestamp, history.Latitude, history.Longitude, history.Accuracy, history.ProviderID, history.UserID).Scan(
		&result.ID,
		&result.Timestamp,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *LocationRepository) DeleteHistory(id, userID int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM locations_history
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("LocationRepository.DeleteHistory: no rows affected")
	}

	return nil
}
