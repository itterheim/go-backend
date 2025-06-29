package services

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/jwt"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	authService *auth.AuthService
	deviceRepo  *repositories.Device
}

func NewDeviceService(deviceRepo *repositories.Device, authService *auth.AuthService) *Device {
	return &Device{
		deviceRepo:  deviceRepo,
		authService: authService,
	}
}

func (s *Device) ListDevices() ([]models.Device, error) {
	return s.deviceRepo.List()
}

func (s *Device) GetDevice(id int64) (*models.Device, error) {
	return s.deviceRepo.GetById(id)
}

func (s *Device) CreateDevice(name string, description string) (*models.Device, error) {
	if len(name) == 0 {
		return nil, errors.New("device.name cannot be empty")
	}

	return s.CreateDevice(name, description)
}

func (s *Device) UpdateDevice(data *models.Device) (*models.Device, error) {
	if len(data.Name) == 0 {
		return nil, errors.New("device.name cannot be empty")
	}

	return s.deviceRepo.Update(data)
}

func (s *Device) DeleteDevice(deviceId int64) error {
	return s.deviceRepo.Delete(deviceId)
}

func (s *Device) CreateToken(deviceId int64, lifespan time.Duration) (string, error) {
	claims := jwt.Claims{
		ID:         deviceId,
		Expiration: time.Now().Add(lifespan),
		JTI:        uuid.New().String(),
		Type:       jwt.DeviceClaim,
	}

	token, err := s.authService.CreateJWTToken(claims)
	if err != nil {
		return "", err
	}

	err = s.deviceRepo.UpdateToken(deviceId, claims.JTI, claims.Expiration)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Device) RevokeToken(deviceId int64) error {
	return s.deviceRepo.RevokeToken(deviceId)
}
