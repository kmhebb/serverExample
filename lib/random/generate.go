package random

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator func() (string, string, error)

// Use of init is generally a smell, but in this case we only want to seed the
// default random number generator, so it's worth keeping it simple.
func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	// PasswordLength is the default length for generated passwords.
	PasswordLength = 8
	// PhraseLength is the number of words to include in a generated passphrase.
	PhraseLength = 3
)

// Bytes returns a sequence of random bytes of length n. It returns an error if
// one occurrs, but this should *never* occur in practice.
// cf. https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SafeBytes returns a sequence of random alphanumeric characters as a slice of
// bytes with the given length.
func SafeBytes(n int) []byte {
	var safeBytes = []byte("1234567890abcdefghijklmnopqrstuvwxyz")
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = safeBytes[rand.Intn(len(safeBytes))]
	}
	return b
}

// String returns a random string of length (s + 2) / 3 * 4.
//
// Since we base64 URL encode the random bytes, we get a padded string always
// longer than the input length.
func String(s int) (string, error) {
	b, err := Bytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Password creates a random, temporary password and returns the plaintext
// password, as well as the bcrypt hashed version.
//
// This function returns an error if one occurs, but this should never really
// happen in practice.
func Password() (string, string, error) {
	newPassBytes, err := Bytes(PasswordLength)
	if err != nil {
		return "", "", errors.Wrap(err, "could not generate new password")
	}

	newPass := base64.URLEncoding.EncodeToString(newPassBytes)
	hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return "", "", errors.Wrap(err, "could not hash password")
	}

	return newPass, string(hash), nil
}

// Words creates a string of n words from a hand-checked list concatenated
// together. The generated string is suitable as a passphrase, used primarily
// for new user passwords and password resets.
func Words(n int) string {
	var ws = make([]string, n)

	for i := range ws {
		j := rand.Int63n(int64(len(wordList)))
		ws[i] = wordList[j]
	}

	return strings.Join(ws, "")
}

// Passphrase return both a generated passphrase and its hash.
func Passphrase() (string, string, error) {
	p := Words(PhraseLength)
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", "", errors.Wrap(err, "could not hash password")
	}

	return p, string(h), nil
}
