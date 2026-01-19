package auth_tokens

import (
	"errors"
	"fmt"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/golang-jwt/jwt/v5"
)

var ErrFailedToSignToken = errors.New("failed to sign token")

type JWTTokenService struct {
	config *auth.Config
}

func NewJWTTokenService(config *auth.Config) *JWTTokenService {
	return &JWTTokenService{
		config: config,
	}
}

// ParseToken validates a JWT token and returns the claims.
func (j *JWTTokenService) ParseToken(tokenStr string) (*auth.JWTClaims, error) {
	if j.config.JwtSecret == "" {
		return nil, auth.ErrJWTNotConfigured
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidSigningMethod
		}

		return []byte(j.config.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", auth.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	// Extract claims
	userID, _ := claims["user_id"].(string)
	sessionID, _ := claims["session_id"].(string)
	exp, _ := claims["exp"].(float64)

	if userID == "" || sessionID == "" {
		return nil, auth.ErrInvalidToken
	}

	return &auth.JWTClaims{
		UserID:    userID,
		SessionID: sessionID,
		ExpiresAt: int64(exp),
	}, nil
}

// GenerateToken creates a new JWT token with the given claims.
func (j *JWTTokenService) GenerateToken(claims *auth.JWTClaims) (string, error) {
	if j.config.JwtSecret == "" {
		return "", auth.ErrJWTNotConfigured
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"session_id": claims.SessionID,
		"exp":        claims.ExpiresAt,
	})

	tokenString, err := token.SignedString([]byte(j.config.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToSignToken, err)
	}

	return tokenString, nil
}
