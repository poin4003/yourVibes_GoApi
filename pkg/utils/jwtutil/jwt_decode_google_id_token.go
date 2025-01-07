package jwtutil

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type GoogleIdTokenClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
	Iat           int64  `json:"iat"`
	jwt.RegisteredClaims
}

func DecodeGoogleIDToken(
	idToken string,
) (*GoogleIdTokenClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(idToken, &GoogleIdTokenClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*GoogleIdTokenClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	if claims.Iss != "https://accounts.google.com" && claims.Iss != "accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", claims.Iss)
	}

	if claims.FamilyName == "" {
		claims.FamilyName = "  "
	}
	if claims.GivenName == "" {
		claims.GivenName = "  "
	}

	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("ID token has expired")
	}

	return claims, nil
}
