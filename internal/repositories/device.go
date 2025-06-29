package repositories

import (
	"backend/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

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
		fmt.Println("Device repository: ", err)
		return nil, nil
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

func (r *Device) UpdateToken(deviceId int64, jti string, expiration time.Time) error {
	commandTag, err := r.db.Exec(context.Background(), `
		UPDATE devices
		SET jti = $2, expiration = $3
		WHERE id = $1
	`, deviceId, jti, expiration)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("UpdateToken: no rows updated")
	}
	return nil
}
