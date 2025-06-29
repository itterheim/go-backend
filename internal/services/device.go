package services

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/jwt"
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
