package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtKey     = "abc"
	refreshKey = "def"

	shortExpDuration = 1 * 24 * time.Hour
	longExpDuration  = 7 * 24 * time.Hour
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func createToken(userID int64, name, key string, dur time.Duration) (string, error) {
	expirationTime := time.Now().Add(dur)

	claims := &Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func CreateAccessToken(userID int64, name string) (string, error) {
	return createToken(userID, name, jwtKey, shortExpDuration)
}

func CreateRefreshToken(userID int64, name string) (string, error) {
	return createToken(userID, name, refreshKey, longExpDuration)
}

func isValidToken(token, key string) (bool, error) {
	var claims Claims
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return false, err
	}

	if !tkn.Valid {
		return false, jwt.ErrSignatureInvalid
	}

	return true, nil
}

func IsValidAccessToken(token string) (bool, error) {
	return isValidToken(token, jwtKey)
}

func IsValidRefreshToken(token string) (bool, error) {
	return isValidToken(token, refreshKey)
}

func UnmarshalPayload(token string) (*Claims, error) {
	tkn := strings.Split(token, ".")
	if len(tkn) != 3 {
		return nil, fmt.Errorf("failed unmarshal jwt payload because wrong format %s", token)
	}

	sDec, err := base64.RawStdEncoding.DecodeString(tkn[1])
	if err != nil {
		return nil, fmt.Errorf("failed decode %s because %s", tkn[1], err.Error())
	}

	var claims Claims
	if err := json.Unmarshal(sDec, &claims); err != nil {
		return nil, fmt.Errorf("failed unmarshal payload %s into struct because %s", tkn[1], err.Error())
	}

	return &claims, nil
}

func Unmarshal(token string, dest interface{}) error {
	tkn := strings.Split(token, ".")
	if len(tkn) != 3 {
		return fmt.Errorf("failed unmarshal jwt payload because wrong format %s", token)
	}

	sDec, err := base64.RawStdEncoding.DecodeString(tkn[1])
	if err != nil {
		return fmt.Errorf("failed decode %s because %s", tkn[1], err.Error())
	}

	if err := json.Unmarshal(sDec, &dest); err != nil {
		return fmt.Errorf("failed unmarshal payload %s into struct because %s", tkn[1], err.Error())
	}

	return nil
}
