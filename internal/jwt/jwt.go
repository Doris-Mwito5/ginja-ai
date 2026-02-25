package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTToken interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type jwtToken struct {
	secretkey string
}

func NewJWTMaker(secretkey string) (JWTToken, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size, must have at least %d characters", minSecretKeySize)
	}

	return &jwtToken{secretkey: secretkey}, nil
}

func (maker *jwtToken) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretkey))
}

func (maker *jwtToken) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}
		return []byte(maker.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errExpiredToken) {
			return nil, errExpiredToken
		}
		return nil, errInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errInvalidToken
	}
	return payload, nil
}
