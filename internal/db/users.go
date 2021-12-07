package db

import (
	"fmt"
	"time"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/pg"
)

var now string = string(time.Now().Format("1/2/2006 15:04"))

func CreateUser(ctx cloud.Context, tx pg.Tx, u *cloud.User) error {
	query := `INSERT INTO users.profile (id, firstname, lastname, email, passhash, mustchange, datecreated, datemodified, lastactivity) VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $7);`

	err := tx.Exec(ctx.Ctx, query, u.ID, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.MustChange, now)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func FindByEmail(ctx cloud.Context, tx pg.Tx, email string) (*cloud.User, error) {
	query := `SELECT CAST(id AS varchar), firstname, lastname, email, passhash, mustchange, lastactivity, datecreated, datemodified FROM users.profile WHERE email = $1;`
	var u cloud.User

	rows, err := tx.Query(ctx.Ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("pg/Tx.UserFindByEmailQuery: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.MustChange, &u.LastActivity, &u.DateModified, &u.DateCreated); err != nil {
			return nil, fmt.Errorf("pg/Tx.UserFindByEmailAssignment: %w", err)
		}
	}
	// There is not an error when the query returns no results, and, oddly, the assignment will not fail. So we will check the result and return it or return a not found error manually.
	if u.Email == email {

		if err = UpdateLastActivity(ctx, tx, u.ID); err != nil {
			// This reports an error in the procedure. Does not interrupt the calling procedure.
			fmt.Println("error updating activity")
		}

		return &u, nil
	}
	return nil, fmt.Errorf("pg/Tx.UserFindByEmailQuery: not found")
}

func FindByID(ctx cloud.Context, tx pg.Tx, id string) (*cloud.User, error) {

	query := `SELECT CAST(id AS varchar), firstname, lastname, email, passhash, mustchange, lastactivity, datecreated, datemodified FROM users.profile WHERE id = $1;`
	var u cloud.User

	rows, err := tx.Query(ctx.Ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("pg/Tx.UserFindByIDQuery: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.MustChange, &u.LastActivity, &u.DateModified, &u.DateCreated); err != nil {
			return nil, fmt.Errorf("pg/Tx.UserFindByIDAssignment: %w", err)
		}
	}

	if err = UpdateLastActivity(ctx, tx, u.ID); err != nil {
		// This reports an error in the procedure. Does not interrupt the calling procedure.
		fmt.Println("error updating activity")
	}

	return &u, nil
}

func UpdateUserRecord(ctx cloud.Context, tx pg.Tx, u *cloud.User) error {
	query := `UPDATE users.profile SET firstname = $2, lastname = $3, email = $4, passhash = $5, mustchange = $6, resettoken = $7, resettokenexpiration = $8, datemodified = $9 WHERE id = $1;`

	err := tx.Exec(ctx.Ctx, query, u.ID, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.MustChange, u.ResetToken, u.ResetTokenExpiration, now)
	if err != nil {
		return fmt.Errorf("pg/Tx.UpdateUserRecord: %w", err)
	}

	if err = UpdateLastActivity(ctx, tx, u.ID); err != nil {
		// This reports an error in the procedure. Does not interrupt the calling procedure.
		fmt.Println("error updating activity")
	}
	return nil
}

func UpdateLastActivity(ctx cloud.Context, tx pg.Tx, uid string) error {
	q := `UPDATE users.profile SET lastactivity = $2 WHERE id = $1`
	err := tx.Exec(ctx.Ctx, q, uid, now)
	if err != nil {
		return err
	}

	return nil
}

func FindByToken(ctx cloud.Context, tx pg.Tx, token string) (*cloud.User, error) {

	query := `SELECT CAST(id AS varchar), firstname, lastname, email, passhash, mustchange, lastactivity, datecreated, datemodified, resettoken, resettokenexpiration FROM users.profile WHERE resettoken = $1;`
	var u cloud.User

	rows, err := tx.Query(ctx.Ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("pg/Tx.UserFindByTokenQuery: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.MustChange, &u.LastActivity, &u.DateCreated, &u.DateModified, &u.ResetToken, &u.ResetTokenExpiration); err != nil {
			return nil, fmt.Errorf("pg/Tx.UserFindByTokenAssignment: %w", err)
		}
	}

	if err = UpdateLastActivity(ctx, tx, u.ID); err != nil {
		// This reports an error in the procedure. Does not interrupt the calling procedure.
		fmt.Println("error updating activity")
	}

	return &u, nil
}

func SaveValidationCode(ctx cloud.Context, tx pg.Tx, email string, code string) error {
	q := `INSERT INTO cache.validation (code, email, datecreated) VALUES ($1, $2, NOW())`
	err := tx.Exec(ctx.Ctx, q, code, email)
	if err != nil {
		return fmt.Errorf("pg/Tx.SaveValidationCode: %w", err)
	}

	return nil
}

func GetValidationCode(ctx cloud.Context, tx pg.Tx, email string) (string, error) {
	var code string

	q := `SELECT code FROM cache.validation WHERE email=$1`
	rows, err := tx.Query(ctx.Ctx, q, email)
	if err != nil {
		return code, fmt.Errorf("pg/Tx.GetValidationCode: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&code); err != nil {
			return code, fmt.Errorf("pg/Tx.GetValidationCodeAssignment: %w", err)
		}
	}

	return code, nil
}

func DeleteValidationCode(ctx cloud.Context, tx pg.Tx, email string) error {
	q := `DELETE FROM cache.validation WHERE email=$1`
	err := tx.Exec(ctx.Ctx, q, email)
	if err != nil {
		return fmt.Errorf("pg/Tx.DeleteValidationCode: %w", err)
	}
	return nil
}

func GetUserList(ctx cloud.Context, tx pg.Tx, listType string) ([]cloud.User, error) {
	var query string
	switch listType {
	case "all":
		query = `SELECT id, email, firstname, lastname, mustchange, datemodified, lastactivity from users.profile`
	default:
		query = `SELECT * from users.profile`
	}

	var users []cloud.User
	rows, err := tx.Query(ctx.Ctx, query)
	if err != nil {
		return []cloud.User{}, fmt.Errorf("pg/Tx.GetUserList query failed: %w", err)
	}

	for rows.Next() {
		var user cloud.User
		if err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MustChange, &user.DateModified, &user.LastActivity); err != nil {
			return []cloud.User{}, fmt.Errorf("pg/Tx.GetUserList Assignment: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
