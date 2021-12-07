package token

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Claims struct {
	// UID is the database ID of the user and identifies a single user of the
	// application.
	UID string
}

// String implements fmt.Stringer.
func (c *Claims) String() string {
	return fmt.Sprintf("<type=Claims,uid=%s>", c.UID)
}

type TokenError struct {
	Cause  error
	Claims jwt.MapClaims
}

func (e *TokenError) Error() string {
	return fmt.Sprintf("%q: %v", e.Claims, e.Cause)
}

type InvalidTokenError struct {
	Cause error
}

func (e *InvalidTokenError) Error() string {
	return e.Cause.Error()
}

var key = []byte(time.Now().Format("20060102150405"))

func SetSigningKey(k string) {
	key = []byte(k)
}

var ErrInvalidSigningMethod = errors.New("invalid signing method")

var (
	// TokenExpiration is set to twenty four hours.
	TokenExpiration = time.Hour * 24

	// OverrideExpiration is set to one day for troubleshooting.
	OverrideExpiration = time.Hour * 4
)

// New creates a new JWT with the given claims, passed as key/value pairs. It
// also sets an expiration time, at present this is 24 hours from issue.
func New(uid string, override string) (string, error) {
	cm := jwt.MapClaims{"sub": uid}
	if override != "" {
		cm["ovr"] = override
		cm["exp"] = int64(time.Now().Add(OverrideExpiration).Unix())
	} else {
		cm["exp"] = int64(time.Now().Add(TokenExpiration).Unix())
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cm)

	tok, err := t.SignedString(key)
	if err != nil {
		return "", &TokenError{Cause: err, Claims: cm}
	}

	return tok, nil
}

// Parse takes a JWT string and returns the claims existing on the JWT. The
// token must be valid and signed with the correct key.
func Parse(h string) (*Claims, error) {
	enc := strings.TrimPrefix(h, "Bearer ")

	// jwt.Parse takes the string representation of the token and a function
	// returning our signing key if the token is valid, or an error otherwise.
	t, err := jwt.Parse(enc, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return key, nil
	})

	if err != nil {
		return nil, &InvalidTokenError{err}
	}

	// The token could potentially still be invalid at this point, so we double
	// check even if err is nil above.
	if !t.Valid {
		return nil, &InvalidTokenError{err}
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &InvalidTokenError{errors.Errorf("invalid claims type: %T", claims)}
	}

	var c Claims
	if c.UID, err = getStringClaim(claims, "sub"); err != nil {
		return nil, &InvalidTokenError{err}
	}

	return &c, nil
}

func getStringClaim(claims jwt.MapClaims, claim string) (string, error) {
	s, ok := claims[claim].(string)
	if !ok {
		return "", errors.Errorf("claim missing or of invalid type: %s", claim)
	}
	return s, nil
}
