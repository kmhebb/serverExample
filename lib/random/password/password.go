package password

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/kmhebb/serverExample/lib/random"
)

var gen random.PasswordGenerator = constant

func Register(g random.PasswordGenerator) {
	gen = g
}

func Generate() (string, string, error) {
	return gen()
}

func constant() (string, string, error) {
	p := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return p, string(hash), nil
}
