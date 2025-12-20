package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func ExtractUserClaimsFromToken(token string) (UserClaims, error) {
	var claims UserClaims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return claims, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return claims, fmt.Errorf("error decoding token: %w", err)
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return claims, fmt.Errorf("error parsing token data: %w", err)
	}

	return claims, nil
}
