package token

import (
	"testing"
	"time"

	"github.com/Franklynoble/bankapp/db/util"
	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/require"
)

//to check JWTMaker
func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.Randomstring(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NotEmpty(t, token)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

//Write Test to Check for Expired Token
func TestExpireToken(t *testing.T) {

	maker, err := NewJWTMaker(util.Randomstring(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// write test for invalidToken
func TestInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.Randomstring(32))

	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
