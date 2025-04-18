package crypto

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func EncryptedPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrapf(err, "error password hash: %v", err)
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

func CheckPassword(plaintext, ciphertext string) bool {
	b, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(b, []byte(plaintext))
	return err == nil
}
