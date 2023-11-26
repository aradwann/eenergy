package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PASETOMaker is PASETO token maker
type PASETOMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPASETOMaker creates a new PasetoMaker
func NewPASETOMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: it must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PASETOMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PASETOMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}
	return token, payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PASETOMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
