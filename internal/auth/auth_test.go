package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	pwd := "testpwdstring"
	hash, err := HashPassword(pwd)
	if err != nil {
		t.Error(err)
	}
	err = CheckPasswordHash(hash, pwd)
	if err != nil {
		t.Error(err)
	}

	wrongPass := "thispasswordiswrong"
	err = CheckPasswordHash(hash, wrongPass)
	if err == nil {
		t.Fatal()
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "dfahjkghfhjgashaghfjkhgajfgl"
	expiresIn := 2 * time.Hour
	testJWT, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}
	res, err := ValidateJWT(testJWT, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	if res != userID {
		t.Fatal()
	}
}
