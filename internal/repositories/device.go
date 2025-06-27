package repositories

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Device struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *Device {
	return &Device{db: db}
}

func (r *Device) GetById(id int64) (*models.Device, error) {
	device := &models.Device{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated, name, description, expiration
		FROM devices
		WHERE username = $1
	`, id).Scan(&device.ID, &device.Created, &device.Updated, &device.Name, &device.Description, &device.Expiration)
	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *Device) Create(name, description string) (*models.Device, error) {
	device := &models.Device{}

	err := r.db.QueryRow(context.Background(), `
		INSERT INTO devices (name, description)
		VALUES ($1, $2)
		RETURNING id, created, updated, name, description
	`, name, description).Scan(&device.ID, &device.Created, &device.Updated, &device.Name, &device.Description)
	if err != nil {
		return nil, err
	}

	return device, err
}

func (r *Device) List() ([]models.Device, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
			id, created, updated, name, description, expiration
		FROM devices
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]models.Device, 0)
	for rows.Next() {
		device := models.Device{}
		err = rows.Scan(&device.ID, &device.Created, &device.Updated, &device.Name, &device.Description, &device.Expiration)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}
