package auth

import (
	"backend/internal/models"
	"backend/internal/repositories"
	jwt1 "backend/pkg/jwt"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   *repositories.User
	tokenRepo  *repositories.Token
	deviceRepo *repositories.Device
	jwtSecret  string
}

func NewAuthService(userRepo *repositories.User, deviceRepo *repositories.Device, tokenRepo *repositories.Token, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
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
	claims := s.createClaims(user.ID, jwt1.UserClaim)
	return s.createJWTTokens(claims)
}

func (s *AuthService) ValidateToken(receivedToken string) (jwt1.Claims, error) {
	// parse token
	token, err := jwt.Parse(receivedToken, func(jwtToken *jwt.Token) (any, error) {
		// validate signing method
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return jwt1.Claims{}, err
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return jwt1.Claims{}, errors.New("invalid token")
	}

	claims, err := s.parseClaims(jwtClaims)
	if err != nil {
		fmt.Println(err)
		return jwt1.Claims{}, err
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

	// rotate refresh token (remove old)
	refresh, access, err := s.createJWTTokens(claims)
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

func (s *AuthService) CreateJWTToken(claims jwt1.Claims) (string, error) {
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.ID,
		"exp":  claims.Expiration.Unix(),
		"jti":  claims.JTI,
		"type": claims.Type,
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) createJWTTokens(claims jwt1.Claims) (refreshToken string, accessToken string, err error) {
	refreshToken, err = s.CreateJWTToken(claims)
	if err != nil {
		return "", "", err
	}

	err = s.tokenRepo.Create(models.Token{
		UserID:     claims.ID,
		Expiration: claims.Expiration,
		JTI:        claims.JTI,
	})
	if err != nil {
		return "", "", err
	}

	// generate new JWT access token
	accessClaims := s.createClaims(claims.ID, claims.Type)
	accessToken, err = s.CreateJWTToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) parseClaims(jwtClaims jwt.MapClaims) (jwt1.Claims, error) {
	claims := jwt1.Claims{}

	// check expiration
	expiration, ok := jwtClaims["exp"].(float64)
	if !ok {
		return jwt1.Claims{}, errors.New("invalid expiration")
	}
	claims.Expiration = time.Unix(int64(expiration), 0)

	// get user ID
	id, ok := jwtClaims["sub"].(float64)
	if !ok {
		return jwt1.Claims{}, errors.New("invalid user ID")
	}
	claims.ID = int64(id)

	// get JTI
	claims.JTI, ok = jwtClaims["jti"].(string)
	if !ok {
		return jwt1.Claims{}, errors.New("invalid JTI")
	}

	// get Type
	claimType, ok := jwtClaims["type"].(string)
	if !ok {
		return jwt1.Claims{}, errors.New("failed to parse claim type")
	}
	claims.Type = jwt1.ClaimType(claimType)
	if claims.Type != jwt1.UserClaim && claims.Type != jwt1.DeviceClaim {
		return jwt1.Claims{}, errors.New("invalid Type")
	}

	return claims, nil
}

func (s *AuthService) createClaims(id int64, claimType jwt1.ClaimType) jwt1.Claims {
	expiration := time.Now().Add(time.Hour * 24)
	jti := uuid.New().String()

	return jwt1.Claims{
		ID:         id,
		JTI:        jti,
		Type:       claimType,
		Expiration: expiration,
	}
}
