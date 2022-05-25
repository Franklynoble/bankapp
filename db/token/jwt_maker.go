package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minsecretKeySize = 32

//JWTMaker is a JSON web Token maker

type JWTMaker struct {
	secretKey string
}

//NewJWTMaker create a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minsecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minsecretKeySize)
	}
	return &JWTMaker{secretKey}, nil

}

// CreateToken creates a new token for a specific username and durati
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	//create new token payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	//this func expect two argument the bit and  the payload
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {

			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}
