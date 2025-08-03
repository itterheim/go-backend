package locations

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaceRepository struct {
	db *pgxpool.Pool
}

func NewPlaceRepository(db *pgxpool.Pool) *PlaceRepository {
	return &PlaceRepository{db}
}

func (r *PlaceRepository) ListPlaces(userID int64) ([]Place, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, name, note, latitude, longitude, radius, created, updated
		FROM locations_places
		ORDER BY name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places := make([]Place, 0)
	for rows.Next() {
		place := Place{}

		err := rows.Scan(&place.ID, &place.Name, &place.Note, &place.Latitude, &place.Longitude, &place.Radius, &place.Created, &place.Updated)
		if err != nil {
			return nil, err
		}

		places = append(places, place)
	}

	return places, nil
}

func (r *PlaceRepository) GetPlace(id, userID int64) (*Place, error) {
	var data Place
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, note, latitude, longitude, radius, created, updated
		FROM locations_places
		WHERE id = $1
	`, id, userID).Scan(&data.ID, &data.Name, &data.Note, &data.Latitude, &data.Longitude, &data.Radius, &data.Created, &data.Updated)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *PlaceRepository) CreatePlace(place *Place) (*Place, error) {
	var result Place
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO locations_places (name, note, latitude, longitude, radius)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, note, latitude, longitude, radius, created, updated
	`, place.Name, place.Note, place.Latitude, place.Longitude, place.Radius).Scan(
		&result.ID, &result.Name, &result.Note, &result.Latitude, &result.Longitude, &result.Radius, &result.Created, &result.Updated,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PlaceRepository) UpdatePlace(place *Place) (*Place, error) {
	var result Place
	err := r.db.QueryRow(context.Background(), `
		UPDATE locations_places
		SET name = $2,
			note = $3,
		    latitude = $4,
			longitude = $5,
			radius = $6
		WHERE id = $1
		RETURNING id, name, note, latitude, longitude, radius, created, updated
	`, place.ID, place.Name, place.Note, place.Latitude, place.Longitude, place.Radius).Scan(
		&result.ID, &result.Name, &result.Note, &result.Latitude, &result.Longitude, &result.Radius, &result.Created, &result.Updated,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PlaceRepository) DeletePlace(id, userID int64) error {
	cmd, err := r.db.Exec(context.Background(), `
		DELETE FROM locations_places
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() != 1 {
		return errors.New("PlaceRepository.DeletePlace: no rows affected")
	}

	return nil
}
