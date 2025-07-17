package service

import (
	"context"
	"errors"
	"go-rest-api-template/internal/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PublicJWTClaims for public endpoints - only contains API key information
type PublicJWTClaims struct {
	ApiKeyID   int    `json:"api_key_id"`
	ApiKeyName string `json:"api_key_name"`
	jwt.RegisteredClaims
}

// PrivateJWTClaims for private endpoints - contains API key and user information
type PrivateJWTClaims struct {
	ApiKeyID   int    `json:"api_key_id"`
	ApiKeyName string `json:"api_key_name"`
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService interface for JWT operations
type JWTService interface {
	GeneratePublicToken(apiKey *entity.ApiKey) (string, error)
	GeneratePrivateToken(apiKey *entity.ApiKey, user *entity.User) (string, error)
	ValidatePublicToken(tokenString string) (*PublicJWTClaims, *entity.ApiKey, error)
	ValidatePrivateToken(tokenString string) (*PrivateJWTClaims, *entity.ApiKey, *entity.User, error)
}

// jwtService implements JWTService
type jwtService struct {
	secretKey              string
	publicTokenExpiration  time.Duration
	privateTokenExpiration time.Duration
	apiKeyService          ApiKeyService
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string, publicExpHours, privateExpHours int, apiKeyService ApiKeyService) JWTService {
	return &jwtService{
		secretKey:              secretKey,
		publicTokenExpiration:  time.Duration(publicExpHours) * time.Hour,
		privateTokenExpiration: time.Duration(privateExpHours) * time.Hour,
		apiKeyService:          apiKeyService,
	}
}

// GeneratePublicToken generates JWT token for public endpoints (API key only)
func (j *jwtService) GeneratePublicToken(apiKey *entity.ApiKey) (string, error) {
	claims := PublicJWTClaims{
		ApiKeyID:   apiKey.ID,
		ApiKeyName: apiKey.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.publicTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-rest-api",
			Subject:   "public-access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GeneratePrivateToken generates JWT token for private endpoints (API key + user)
func (j *jwtService) GeneratePrivateToken(apiKey *entity.ApiKey, user *entity.User) (string, error) {
	claims := PrivateJWTClaims{
		ApiKeyID:   apiKey.ID,
		ApiKeyName: apiKey.Name,
		UserID:     user.ID,
		Username:   user.Username,
		Email:      user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.privateTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-rest-api",
			Subject:   "private-access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidatePublicToken validates public JWT token and returns API key entity
func (j *jwtService) ValidatePublicToken(tokenString string) (*PublicJWTClaims, *entity.ApiKey, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PublicJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*PublicJWTClaims)
	if !ok {
		return nil, nil, errors.New("invalid token claims")
	}

	// Get API key entity from service
	ctx := context.Background()
	apiKey, err := j.apiKeyService.GetApiKeyByID(ctx, claims.ApiKeyID)
	if err != nil {
		return nil, nil, errors.New("api key not found")
	}

	return claims, apiKey, nil
}

// ValidatePrivateToken validates private JWT token and returns API key + user entities
func (j *jwtService) ValidatePrivateToken(tokenString string) (*PrivateJWTClaims, *entity.ApiKey, *entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PrivateJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	if !token.Valid {
		return nil, nil, nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*PrivateJWTClaims)
	if !ok {
		return nil, nil, nil, errors.New("invalid token claims")
	}

	// Validate that this is actually a private token by checking user fields
	if claims.UserID == 0 || claims.Username == "" || claims.Email == "" {
		return nil, nil, nil, errors.New("token is not a private token")
	}

	// Get API key entity from service
	ctx := context.Background()
	apiKey, err := j.apiKeyService.GetApiKeyByID(ctx, claims.ApiKeyID)
	if err != nil {
		return nil, nil, nil, errors.New("api key not found")
	}

	// Create user entity from claims (since we have the user data in the token)
	user := &entity.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}

	return claims, apiKey, user, nil
}
