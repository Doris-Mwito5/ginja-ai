package jwt

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

var (
	errExpiredToken = errors.New("token has expired")
	errInvalidToken = errors.New("token is invalid")
)


type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID := uuid.NewRandom()
	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

//validate the payload
func (payload *Payload) Valid() error {
	if time.Now().After((payload.ExpiresAt)) {
		return errExpiredToken
	}
	return nil
}

