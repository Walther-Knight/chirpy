package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pwd := []byte(password)
	hashedPass, errHash := bcrypt.GenerateFromPassword(pwd, 3)
	if errHash != nil {
		log.Printf("error hashing password: %v", errHash)
		return "", errHash
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		//expiresIn defined in api.UserLogin()
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	secretKey := []byte(tokenSecret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("error creating auth token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("authentication error decoding token: %v", err)
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	tokenString, found := strings.CutPrefix(headers.Get("Authorization"), "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if !found {
		log.Printf("invalid or missing token string")
		return "", errors.New("invalid or missing token string")
	}
	return tokenString, nil

}

func MakeRefreshToken() (string, error) {
	tokenSeed := make([]byte, 32)
	_, err := rand.Read(tokenSeed)
	if err != nil {
		log.Printf("error seeding refresh token: %v", err)
	}
	return hex.EncodeToString(tokenSeed), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiString, found := strings.CutPrefix(headers.Get("Authorization"), "ApiKey ")
	apiString = strings.TrimSpace(apiString)
	if !found {
		log.Printf("missing apikey string")
		return "", errors.New("missing ApiKey string")
	}
	return apiString, nil
}
