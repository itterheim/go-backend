package core

import (
	"backend/internal/config"
	pkgjwt "backend/pkg/jwt"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo     *UserRepository
	tokenRepo    *TokenRepository
	providerRepo *ProviderRepository
	config       *config.AuthConfig
}

func NewAuthService(userRepo *UserRepository, providerRepo *ProviderRepository, tokenRepo *TokenRepository, config *config.AuthConfig) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		providerRepo: providerRepo,
		tokenRepo:    tokenRepo,
		config:       config,
	}
}

// Login = find user, verify password, create JWT refresh token
func (s *AuthService) Login(username, password string) (string, string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", "", err
	}

	if !user.CheckPassword(password) {
		return "", "", errors.New("invalid password")
	}

	// create new tokens
	accessClaims := s.createClaims(user.ID, pkgjwt.UserClaim, s.config.AccessExpiration)
	return s.createJWTTokens(accessClaims)
}

func (s *AuthService) ValidateToken(receivedToken string) (pkgjwt.Claims, error) {
	// parse token
	token, err := jwt.Parse(receivedToken, func(jwtToken *jwt.Token) (any, error) {
		// validate signing method
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return pkgjwt.Claims{}, err
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return pkgjwt.Claims{}, errors.New("invalid token")
	}

	claims, err := s.parseClaims(jwtClaims)
	if err != nil {
		fmt.Println(err)
		return pkgjwt.Claims{}, err
	}

	return claims, nil
}

func (s *AuthService) ValidateRefreshToken(token string) (string, string, error) {
	// validate JWT token
	claims, err := s.ValidateToken(token)
	if err != nil {
		return "", "", err
	}

	// find stored token
	t, err := s.tokenRepo.Get(claims.JTI)
	if err != nil {
		return "", "", err
	}

	if t.Blocked {
		return "", "", errors.New("token blocked")
	}

	user, err := s.userRepo.GetUser(claims.UserID)
	if err != nil {
		return "", "", errors.New("ValidateRefreshToken: failed to retrieve user")
	}

	accessClaims := s.createClaims(user.ID, pkgjwt.UserClaim, s.config.AccessExpiration)

	// rotate refresh token (remove old)
	refresh, access, err := s.createJWTTokens(accessClaims)
	if err != nil {
		return "", "", err
	}

	// remove old token
	err = s.tokenRepo.Delete(claims.JTI)
	if err != nil {
		return "", "", err
	}

	return refresh, access, nil
}

func (s *AuthService) CreateJWTToken(claims pkgjwt.Claims) (string, error) {
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.UserID,
		"pid":  claims.ProviderID,
		"exp":  claims.Expiration.Unix(),
		"jti":  claims.JTI,
		"type": claims.Type,
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) createJWTTokens(claims pkgjwt.Claims) (refreshToken string, accessToken string, err error) {
	refreshClaims := s.createClaims(claims.UserID, claims.Type, s.config.RefreshExpiration)
	refreshToken, err = s.CreateJWTToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	err = s.tokenRepo.Create(Token{
		UserID:     refreshClaims.UserID,
		Expiration: refreshClaims.Expiration,
		JTI:        refreshClaims.JTI,
	})
	if err != nil {
		return "", "", err
	}

	// generate new JWT access token
	accessClaims := s.createClaims(claims.UserID, claims.Type, s.config.AccessExpiration)
	accessToken, err = s.CreateJWTToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) parseClaims(jwtClaims jwt.MapClaims) (pkgjwt.Claims, error) {
	claims := pkgjwt.Claims{}

	// check expiration
	expiration, ok := jwtClaims["exp"].(float64)
	if !ok {
		return pkgjwt.Claims{}, errors.New("invalid expiration")
	}
	claims.Expiration = time.Unix(int64(expiration), 0)

	// get user ID
	id, ok := jwtClaims["sub"].(float64)
	if !ok {
		return pkgjwt.Claims{}, errors.New("invalid user ID")
	}
	claims.UserID = int64(id)

	// get JTI
	claims.JTI, ok = jwtClaims["jti"].(string)
	if !ok {
		return pkgjwt.Claims{}, errors.New("invalid JTI")
	}

	// get Type
	claimType, ok := jwtClaims["type"].(string)
	if !ok {
		return pkgjwt.Claims{}, errors.New("failed to parse claim type")
	}
	claims.Type = pkgjwt.ClaimType(claimType)
	if claims.Type != pkgjwt.UserClaim && claims.Type != pkgjwt.ProviderClaim {
		return pkgjwt.Claims{}, errors.New("invalid Type")
	}

	if claims.Type == pkgjwt.ProviderClaim {
		providerId, ok := jwtClaims["pid"].(float64)
		if !ok {
			return pkgjwt.Claims{}, errors.New("invalid provider ID")
		}
		parsed := int64(providerId)
		claims.ProviderID = &parsed
	}

	return claims, nil
}

func (s *AuthService) createClaims(id int64, claimType pkgjwt.ClaimType, expiresInMinutes time.Duration) pkgjwt.Claims {
	expiration := time.Now().Add(time.Minute * expiresInMinutes)
	jti := uuid.New().String()

	return pkgjwt.Claims{
		UserID:     id,
		Type:       claimType,
		JTI:        jti,
		Expiration: expiration,
	}
}
