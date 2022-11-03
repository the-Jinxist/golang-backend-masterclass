package token

import (
	"errors"
	"time"

	uuid "github.com/google/uuid"
)

var (
	ErrorExpiredToken = errors.New("token has expired")
	ErrInvalidToken   = errors.New("invalid token")
)

//Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

//NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

//This method is needed for the claims part of jwt.NewTokenWithClaims. To pass your payload as a claims file, it needs
//to override this Valid() method
func (paylod *Payload) Valid() error {
	if time.Now().After(paylod.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
