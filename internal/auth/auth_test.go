package auth

import (
	"testing"
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
