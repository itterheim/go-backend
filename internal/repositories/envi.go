package repositories

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Envi struct {
	db *pgxpool.Pool
}

func NewEnviRepository(db *pgxpool.Pool) *Envi {
	return &Envi{db: db}
}

func (r *Envi) GetById(id int64) (*models.Device, error) {
	device := &models.Device{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated, name
		FROM envi_devices
		WHERE id = $1
	`, id).Scan(&device.ID, &device.Created, &device.Updated, &device.Name)
	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return device, nil
}
