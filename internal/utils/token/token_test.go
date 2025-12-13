package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testSecretKey = "secretKey"
	testUser      = "TestUser"
	tokenDuration = 12 * time.Hour
	invalidToken  = "invalid.Token"
)

func TestCreateToken(t *testing.T) {
	JWTManager := NewJWTManager(testSecretKey, tokenDuration)
	testToken, err := JWTManager.CreateToken(testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, testToken)

	token, err := jwt.ParseWithClaims(testToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(testSecretKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*UserClaims)
	assert.True(t, ok)
	assert.Equal(t, testUser, claims.Login)
}

func TestVerifyToken(t *testing.T) {
	JWTManager := NewJWTManager(testSecretKey, tokenDuration)
	testToken, err := JWTManager.CreateToken(testUser)
	assert.NoError(t, err)

	userClaims, err := JWTManager.VerifyToken(testToken)
	assert.NoError(t, err)
	assert.Equal(t, testUser, userClaims.Login)

	_, err = JWTManager.VerifyToken(invalidToken)
	assert.Error(t, err)

	_, err = JWTManager.VerifyToken("")
	assert.Error(t, err)
}

func TestExtractUserFromToken(t *testing.T) {
	JWTManager := NewJWTManager(testSecretKey, tokenDuration)
	testToken, err := JWTManager.CreateToken(testUser)
	assert.NoError(t, err)

	userClaims, err := ExtractUserClaimsFromToken(testToken)
	assert.NoError(t, err)
	assert.Equal(t, testUser, userClaims.Login)

	_, err = ExtractUserClaimsFromToken(invalidToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")

	_, err = ExtractUserClaimsFromToken("")
	assert.Error(t, err)
}
