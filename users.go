package cloud

type User struct {
	ID                   string `json:"id" db:"id"`
	Email                string `json:"email" db:"email"`
	FirstName            string `json:"firstName" db:"firstname"`
	LastName             string `json:"lastName" db:"lastname"`
	PasswordHash         string `json:"-" db:"passhash"`
	MustChange           bool   `json:"mustChange" db:"mustchange"`
	DateCreated          string `json:"dateCreated" db:"datecreated"`
	DateModified         string `json:"dateModified" db:"datemodified"`
	LastActivity         string `json:"lastActivity" db:"lastactivity"`
	ResetToken           string `json:"-" db:"resettoken"`
	ResetTokenExpiration string `json:"-" db:"resettokenexpiration"`
}
