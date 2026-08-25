package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTPayloadClaims struct {
	Payload
	jwt.RegisteredClaims
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func NewJWTPayloadClaims(payload *Payload) *JWTPayloadClaims {
	return &JWTPayloadClaims{
		Payload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(payload.IssuedAt),
			Issuer:    "simplebank",
			Subject:   payload.Username,
			ID:        payload.ID.String(),
			Audience:  []string{"clients"},
		},
	}
}

// Create creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	payloadClaims := NewJWTPayloadClaims(payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadClaims)
	signedString, err := jwtToken.SignedString([]byte(maker.secretKey))
	return signedString, nil
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (any, error) {
		_, ok := (token.Method).(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayloadClaims{}, keyfunc)
	if err != nil {
		if err == jwt.ErrTokenInvalidClaims {
			return nil, ErrInvalidToken
		}
		return nil, ErrExpiredToken
	}

	jwtPayloadClaims, ok := jwtToken.Claims.(*JWTPayloadClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        jwtPayloadClaims.Payload.ID,
		Username:  jwtPayloadClaims.Payload.Username,
		IssuedAt:  jwtPayloadClaims.Payload.IssuedAt,
		ExpiredAt: jwtPayloadClaims.Payload.ExpiredAt,
	}
	return payload, nil
}
