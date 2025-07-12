package core

import (
	"backend/pkg/jwt"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ProviderService struct {
	authService  *AuthService
	providerRepo *ProviderRepository
}

func NewProviderService(providerRepo *ProviderRepository, authService *AuthService) *ProviderService {
	return &ProviderService{
		providerRepo: providerRepo,
		authService:  authService,
	}
}

func (s *ProviderService) ListProviders() ([]Provider, error) {
	return s.providerRepo.List()
}

func (s *ProviderService) GetProvider(id int64) (*Provider, error) {
	return s.providerRepo.GetById(id)
}

func (s *ProviderService) CreateProvider(name string, description string) (*Provider, error) {
	if len(name) == 0 {
		return nil, errors.New("provider.name cannot be empty")
	}

	return s.CreateProvider(name, description)
}

func (s *ProviderService) UpdateProvider(data *Provider) (*Provider, error) {
	if len(data.Name) == 0 {
		return nil, errors.New("provider.name cannot be empty")
	}

	return s.providerRepo.Update(data)
}

func (s *ProviderService) DeleteProvider(providerId int64) error {
	return s.providerRepo.Delete(providerId)
}

func (s *ProviderService) CreateToken(providerId int64, lifespan time.Duration) (string, error) {
	claims := jwt.Claims{
		UserID:     providerId,
		Expiration: time.Now().Add(lifespan),
		JTI:        uuid.New().String(),
		Type:       jwt.ProviderClaim,
	}

	token, err := s.authService.CreateJWTToken(claims)
	if err != nil {
		return "", err
	}

	err = s.providerRepo.UpdateToken(providerId, claims.JTI, claims.Expiration)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *ProviderService) RevokeToken(providerId int64) error {
	return s.providerRepo.RevokeToken(providerId)
}
