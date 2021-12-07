package users

import (
	"time"

	cloud "github.com/kmhebb/serverExample"
	password "github.com/kmhebb/serverExample/lib/random/password"
)

func NewUser(email string, firstName string, lastName string) (*cloud.User, error) {
	var user cloud.User

	user.DateCreated = string(time.Now().Format("1/2/2006"))
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName

	return &user, nil
}

func ResetPassword(u cloud.User) (string, error) {
	password, hash, err := password.Generate()
	if err != nil {
		return "", err
	}

	u.MustChange = true
	u.PasswordHash = hash
	u.ResetToken = ""
	u.ResetTokenExpiration = ""

	return password, nil
}
