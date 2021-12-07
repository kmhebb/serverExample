package service

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/internal/db"
	"github.com/kmhebb/serverExample/lib/email"
	"github.com/kmhebb/serverExample/lib/random"
	"github.com/kmhebb/serverExample/lib/random/password"
	"github.com/kmhebb/serverExample/lib/token"
	"github.com/kmhebb/serverExample/log"
	"github.com/kmhebb/serverExample/pg"
	"github.com/pborman/uuid"
)

type UserService struct {
	DB Database
	L  log.Logger
	Em email.Service
}

type GetUserRequest struct {
	ID string `json:"id"`
}

type GetUserResponse struct {
	User *cloud.User `json:"user"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type PutUserRequest struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

type UserRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token string `json:"token"`
}

type ConfirmValidationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type CreateNewUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type NewUserResponse struct {
	User *cloud.User `json:"user"`
}

type ListUsersRequest struct {
	ListType string `json:"listType"`
}

type ListUsersResponse struct {
	Users []cloud.User
}

func (svc UserService) CreateNewUser(ctx cloud.Context, req CreateNewUserRequest) (*NewUserResponse, *cloud.Error) {
	// We will need to come up with a way to authenticate these requests. Will depend on if it will be user driven or admin driven.

	if req.Email == "" {
		return &NewUserResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "email is required",
			//Cause: err,
		}) //fmt.Errorf("email is required")
	}

	user, err := svc.FindOrCreate(ctx, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return &NewUserResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service failed to create new user",
			Cause:   err,
		}) //fmt.Errorf("failed to create new user: %s", err)
	}

	return &NewUserResponse{User: user}, nil
}

func (svc UserService) FindOrCreate(ctx cloud.Context, email string, firstName string, lastName string) (*cloud.User, *cloud.Error) {
	if email == "" {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "email is required",
			//Cause: err,
		}) //fmt.Errorf("email is required")
	}

	email = strings.ToLower(email)

	var existing *cloud.User
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		existing, dbErr = db.FindByEmail(ctx, tx, email)
		return dbErr
	})
	// Strange quirk here. If the find succeeds, ie err == nil, we need to return the existing profile.
	// If err != nil, the find failed, we have to continue with then new user creation.
	// This complicates error reporting if there is truly a DB error. But that will be another day.
	if err == nil {
		return existing, nil
	}

	newUser, err := NewUser(email, firstName, lastName)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service failed to initialize new user",
			Cause:   err,
		}) //err
	}
	passphrase, err := SetTempPassword(newUser)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service failed to set a temp password",
			Cause:   err,
		}) //fmt.Errorf("failed to set passphrase")
	}
	newUser.ID = uuid.New()
	newUser.DateCreated = string(time.Now().Format("1/2/2006"))
	newUser.DateModified = string(time.Now().Format("1/2/2006"))

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.CreateUser(ctx, tx, newUser)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service create user db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.CreateUser.RunInTransaction failed: %w", err)
	}

	svc.L.Info(ctx.Ctx, "Created new user", log.Fields{"email": email})

	svc.Em.NewUserAsync(ctx, svc.Name(newUser), email, passphrase)

	return newUser, nil
}

func (svc UserService) Get(ctx cloud.Context, req GetUserRequest) (resp GetUserResponse, e *cloud.Error) {
	// Only users allowed to access this data should request it.
	// Well, we did check in the handler that there was a token, but now we will validate that the user profile is valid.
	accessError := svc.ValidateUserAuth(ctx)
	if accessError != nil {
		return GetUserResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   accessError,
		}) //fmt.Errorf("user profile not validated: %w", accessError)
	}

	if req.ID == "" {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "id is required",
			//Cause: err,
		}) //fmt.Errorf("id is required")
	}

	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		resp.User, dbErr = db.FindByID(ctx, tx, req.ID)
		return dbErr
	})
	if err != nil {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service find by id db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.GetUser.RunInTransaction failed: %w", err)
	}

	return resp, nil
}

func (svc UserService) Login(ctx cloud.Context, req LoginUserRequest) (resp LoginUserResponse, e *cloud.Error) {
	if req.Email == "" {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "email is required",
			//Cause:   err,
		}) //fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "password is required",
			//Cause:   err,
		}) //fmt.Errorf("password is required")
	}

	var u *cloud.User
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		u, dbErr = db.FindByEmail(ctx, tx, req.Email)
		return dbErr
	})
	if err != nil {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service find by email db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.Login.RunInTransaction failed: %w", err)
	}

	ctx.UserKey = u.ID

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			svc.L.Info(ctx.Ctx, "Failed to verify password", log.Fields{"email": req.Email, "err": err})
			return resp, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindBadRequest,
				Message: "invalid password - password did not match",
				Cause:   err,
			}) //err
		}

		// We are not tracking failed logins, but we could. I will leave this here.
		// err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		// 	return db.UpdateLastFailedLogin(ctx, tx, u)
		// })
		// if err != nil {
		// 	svc.L.Info(ctx.Ctx, "Failed to update last failed login", log.Fields{"email": req.Email, "err": err})
		// }

		svc.L.Info(ctx.Ctx, "Login Failed", log.Fields{"email": req.Email})
		// TODO: install auditor and log this.

		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "error comparing password",
			Cause:   err,
		}) //err
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.UpdateLastActivity(ctx, tx, u.ID)
	})
	if err != nil {
		svc.L.Info(ctx.Ctx, "Failed to update last activity", log.Fields{"email": req.Email, "err": err})
	}

	svc.L.Info(ctx.Ctx, "Login succeeded", log.Fields{"user": req.Email})
	// audit entry here too.

	resp.Token, err = token.New(string(u.ID), "")
	if err != nil {
		return resp, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service failed to issue valid token",
			Cause:   err,
		}) //err
	}

	return resp, nil
}

func (svc UserService) Put(ctx cloud.Context, req PutUserRequest) (interface{}, *cloud.Error) {
	// Only users allowed to access this data should request it.
	// Well, we did check in the handler that there was a token, but now we will validate that the user profile is valid.
	accessError := svc.ValidateUserAuth(ctx)
	if accessError != nil {
		return GetUserResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   accessError,
		}) //fmt.Errorf("user profile not validated: %w", accessError)
	}
	// This endpoint will allow users to update user account profile values. Email and password are the only tricky ones, really.
	if req.ID == "" {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "id is required",
			//Cause: err,
		}) //fmt.Errorf("id is required")
	}

	// First thing we will do is pull up the user profile. Then we will figure out what the user wants to change and then commit those things to the db.
	var u *cloud.User
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		u, dbErr = db.FindByID(ctx, tx, req.ID)
		return dbErr
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service find by id db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.Put.FindByID.RunInTransaction failed: %w", err)
	}
	ctx.UserKey = u.ID

	// Lets check the email first.
	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		// if the profile email and request email dont match, they want to change it.
		// but we are going to check and make sure it is not already in use on the platform. We do not want duplicates.
		if u.Email != req.Email {
			var existing *cloud.User
			err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
				var dbErr error
				existing, dbErr = db.FindByEmail(ctx, tx, req.Email)
				return dbErr
			})
			// Strange quirk here. If the find succeeds, ie err == nil, that means it is in use in the db.
			// So we will need to compare the id returned to the id submitted. If they are the same, its the same profile.
			// If they are different, it means the user wants to change the email to an email that is already in use in the system. We will not allow this.
			if err == nil {
				if req.ID != existing.ID {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "email already in use on this platform, choose a different email",
						Cause:   err,
					}) //fmt.Errorf("email already in use on this platform, choose a different email")
				}
				// If there is an error, ie err != nil, that is fine, the user can change the email to that value.
			}
			u.Email = req.Email
		}
	}

	// Password is a property that we will commonly want to change, but is optional for this procedure - the user could change other properties and not PW.
	// But, in cases where this call is to change the PW we want to check the old pw hash first, to make sure the user requesting this can change it.
	// Then, if that succeeds, we will need to take a hash of the value and save that to the db.
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindBadRequest,
				Message: "old password is required in order to change password",
				//Cause:   err,
			}) //fmt.Errorf("old password is required to change password")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.OldPassword)); err != nil {
			if err != bcrypt.ErrMismatchedHashAndPassword {
				svc.L.Info(ctx.Ctx, "Failed to verify password", log.Fields{"email": req.Email, "err": err})
				return nil, cloud.NewError(cloud.ErrOpts{
					Kind:    cloud.ErrKindInternal,
					Message: "failed to verify old password - try again",
					Cause:   err,
				}) //fmt.Errorf("failed to verify old password - try again")
			}
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindBadRequest,
				Message: "old password was not correct - try again",
				Cause:   err,
			}) //fmt.Errorf("old password was not correct - try again")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindInternal,
				Message: "there was an error hashing the new password",
				Cause:   err,
			}) //fmt.Errorf("there was an error hashing password: %w", err)
		}
		u.PasswordHash = string(hash)
		u.MustChange = false
	}

	// The name properties are trivial. We will let the user do whatever they want with them. If they are populated, we will save those values.
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.UpdateUserRecord(ctx, tx, u)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service update user db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.PutUserRecord.RunInTransaction failed: %w", err)
	}

	svc.L.Info(ctx.Ctx, "Updated user", log.Fields{"id": req.ID, "email": req.Email})
	//  TODO: s.auditor.Log(ctx,
	// 	"Update User Succeeded",
	// 	u.Email,
	// 	"user",
	// 	"update")

	return nil, nil
}

func (svc UserService) RequestPasswordReset(ctx cloud.Context, req UserRequest) (interface{}, *cloud.Error) {
	if req.Email == "" {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "email is required",
			//Cause: err,
		}) //fmt.Errorf("email is required")
	}

	req.Email = strings.ToLower(req.Email)
	var u *cloud.User
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		u, dbErr = db.FindByEmail(ctx, tx, req.Email)
		return dbErr
	})
	if err != nil {
		return "", cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service find by email db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.RequestPwdReset.RunInTransaction failed: %w", err)
	}

	token, _ := random.String(cloud.DefaultTokenLength)
	u.ResetToken = token
	u.ResetTokenExpiration = string(time.Now().Add(24 * time.Hour).Format("1/2/2006 15:04"))

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.UpdateUserRecord(ctx, tx, u)
	})

	if err != nil {
		return ResetFailed, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service update user db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.ResetPassword.RunInTransaction failed: %w", err)
	}

	svc.L.Info(ctx.Ctx, "User requested a password reset, token set", log.Fields{"email": u.Email})

	// s.auditor.Log(ctx,
	// 	"Password Reset Requested",
	// 	"password reset initiated",
	// 	"user",
	// 	"reset_password",
	// )

	svc.Em.ResetPasswordAsync(ctx, svc.Name(u), u.Email, u.ResetToken)

	return nil, nil
}

func (svc UserService) ResetPassword(ctx cloud.Context, req ResetPasswordRequest) (interface{}, *cloud.Error) {
	// we have already checked that the confirmation token is present.

	var u *cloud.User
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		u, dbErr = db.FindByToken(ctx, tx, ctx.ConfirmationToken)
		return dbErr
	})
	if err != nil {
		return ResetFailed, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service find reset token db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.ResetPassword.FindByToken.RunInTransaction failed: %w", err)
	}

	if u.Email == "" {
		return ResetFailed, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "reset token invalid, user not found",
		})
	}

	// ok. the confirmation token was valid and we have a userid, we will now generate a password, email it to the user
	// and then we will save the hash of that password to the database for the user to login with.
	ctx.UserKey = u.ID

	password, hash, err := password.Generate()
	if err != nil {
		return ResetFailed, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "failed to generate temp password",
			//Cause:   err,
		}) //err
	}

	u.MustChange = true
	u.PasswordHash = hash
	u.ResetToken = ""
	u.ResetTokenExpiration = ""

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.UpdateUserRecord(ctx, tx, u)
	})

	if err != nil {
		return ResetFailed, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service update user record db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.ResetPassword.UpdateUserRecord.RunInTransaction failed: %w", err)
	}

	svc.L.Info(ctx.Ctx, "Reset password for user", log.Fields{"email": u.Email})

	svc.Em.NewPasswordAsync(ctx, svc.Name(u), u.Email, password)

	return resetSucceeded, nil
}

func (svc UserService) ListUsers(ctx cloud.Context, req ListUsersRequest) (interface{}, *cloud.Error) {
	// Only users allowed to access this data should request it.
	// Well, we did check in the handler that there was a token, but now we will validate that the user profile is valid.
	accessError := svc.ValidateUserAuth(ctx)
	if accessError != nil {
		return GetUserResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   accessError,
		}) //fmt.Errorf("user profile not validated: %w", accessError)
	}

	var resp ListUsersResponse
	var err error
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		resp.Users, err = db.GetUserList(ctx, tx, req.ListType)
		if err != nil {
			return fmt.Errorf("service/db.GetGridBatchList failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return GridBatchListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service list users db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.ListUsers.RunInTransaction failed: %w", err)
	}

	return resp, nil
}

func (svc UserService) ValidateUserAuth(ctx cloud.Context) error {
	var u *cloud.User
	var err error
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		u, err = db.FindByID(ctx, tx, ctx.UserKey)
		if err != nil {
			return fmt.Errorf("userservice.validateuserauth.findbyid failed: %w", err)
		}
		return nil
	})
	if u.ID == "" {
		return fmt.Errorf("user not valid")
	}
	return nil
}

// func (svc UserService) Find(ctx cloud.Context, req UserRequest) (GetUserResponse, *cloud.Error) {
// 	var err error
// 	resp := GetUserResponse{}
// 	resp.User, err = svc.Es.Find(ctx, ctx.CID, ctx.UserKey)
// 	if err != nil {
// 		return resp, fmt.Errorf("failed to retreive employee information: %s", err)
// 	}

// 	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
// 		var dbErr error
// 		resp.User, dbErr = db.FindByID(ctx, tx, ctx.UserKey)
// 		return dbErr
// 	})
// 	if err != nil {
// 		return resp, fmt.Errorf("service/UserService.GetUser.RunInTransaction failed: %w", err)
// 	}

// 	return resp, nil
// }

func (svc UserService) ValidateEmail(ctx cloud.Context, req UserRequest) (interface{}, error) {
	req.Email = strings.ToLower(req.Email)
	code := random.GenerateCode(6)

	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		dbErr := db.SaveValidationCode(ctx, tx, req.Email, code)
		return dbErr
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service save validation code transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.SaveValidation.RunInTransaction failed: %w", err)
	}

	svc.Em.ValidateEmailAsync(ctx, req.Email, code)

	return nil, nil
}

func (svc UserService) ConfirmValidation(ctx cloud.Context, req ConfirmValidationRequest) (interface{}, error) {
	var gotCode string
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		gotCode, dbErr = db.GetValidationCode(ctx, tx, req.Email)
		return dbErr
	})
	if err != nil {
		if err.Error() == "no validation code for email" {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindBadRequest,
				Message: "no validation code for email",
				Cause:   err,
			}) //fmt.Errorf("no validation code for email: %s", req.Email)
		}
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "user service get validation code db transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/UserService.GetValidation.RunInTransaction failed: %w", err)
	}

	if req.Code == gotCode {
		err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
			dbErr := db.DeleteValidationCode(ctx, tx, req.Email)
			return dbErr
		})
		if err != nil {
			// we are not  going to throw off the whole process by returning an error here. log it, see if it happens much.
			svc.L.Error(ctx.Ctx, err, "validation code not deleted", log.Fields{"code: ": req.Code, "email: ": req.Email})
		}

		return nil, nil
	}

	return nil, fmt.Errorf("code mismatch")
}

func (svc UserService) UpdateLastActivity(ctx cloud.Context, req GetUserRequest) error {
	if req.ID == "" {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "id is required",
			//Cause: err,
		}) //fmt.Errorf("id is required")
	}

	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		dbErr := db.UpdateLastActivity(ctx, tx, ctx.UserKey)
		return dbErr
	})
	if err != nil {
		return fmt.Errorf("service/UserService.UpdateLastActivity.RunInTransaction failed: %w", err)
	}

	return nil
}

func (svc UserService) Name(u *cloud.User) string {
	return strings.Join([]string{u.FirstName, u.LastName}, " ")
}

func NewUser(email, fn, ln string) (*cloud.User, error) {
	var fieldErrors []*FieldError
	if email == "" {
		fieldErrors = append(fieldErrors, NewFieldError("Email", "Required"))
	}
	if fn == "" {
		fieldErrors = append(fieldErrors, NewFieldError("First Name", "Required)"))
	}
	if ln == "" {
		fieldErrors = append(fieldErrors, NewFieldError("Last Name", "Required)"))
	}
	if len(fieldErrors) > 0 {
		return nil, fmt.Errorf("%s", fieldErrors)
	}

	return &cloud.User{
		Email:     strings.ToLower(email),
		FirstName: fn,
		LastName:  ln,
	}, nil
}

func NewFieldError(name string, errors ...string) *FieldError {
	if len(errors) < 1 {
		return nil
	}
	fe := &FieldError{Name: name}
	for _, err := range errors {
		if err != "" {
			fe.Errors = append(fe.Errors, err)
		}
	}
	return fe
}

// FieldError represents any kind of validation error for a specific field. Name
// should be a user-friendly representation of the field name and Errors should
// be client-safe messages describing which validation requirement failed. All
// field validation should be performed in one go and all errors returned as a
// single FieldError to prevent a user having to resend a request multiple times
// to create a valid domain object.
//
// The string value of a FieldError looks similar to:
//
//   Name: Required, Length must be greater than 10
type FieldError struct {
	Name   string
	Errors []string
}

func (e *FieldError) Error() string {
	if len(e.Errors) < 1 {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Name, strings.Join(e.Errors, ", "))
}

func SetTempPassword(u *cloud.User) (string, error) {
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

// HTML Templates

var ResetFailed = `
<!doctype html>

<html>
    <head>
        <title>
            Password Reset Failed
        </title>
        <style>
            * {
                margin: 0;
                padding: 0;
            }

            html {
                height: 100%;
            }
            
            body {
                height: 100%;
            }

            main {
                display: flex;
                flex-direction: column;
                height: 100%;
            }

            nav {
                background: #000;
                height: 50px;
                padding: 10px;
                display: flex;
                align-items: center;
                border-top: 5px solid #ffd600;
            }

            .logo {
                height: 50px;
            }

            h1 {
                color: #ffd600;
            }

            .content {
                background: #fff;
                flex: 1;
                flex-direction: column;
                height: 100%;
                display: flex;
                justify-content: center;
                align-items: center;
            }

            p {
                font-family: 'Lato', sans-serif;
                color: #333;
                padding: 10px;
                font-size: 20px;
            }
        </style>
        <link href="https://fonts.googleapis.com/css?family=Lato:400" rel="stylesheet" type="text/css">
    </head>
    <body>
        <main>
            <nav>
                <a href="https://npandl.com">
                    <img class="logo" src="data:image/png;base64," />
                </a>
            </nav>
            <div class="content">
                <p>There was an error resetting your password.</p>
                <p>Please try again or contact support if the issue persists at:</p>
                <p></p>
                <p><a href="mailto:">mail@npandl.com</a></p>
                <p>-or-</p>
                <p><a href="tel:">123-456-7890</a></p>
            </div>
        </main>
    </body>
</html>
`

var resetSucceeded = `
<!doctype html>

<html>
    <head>
        <title>
            Password Reset Successful
        </title>
        <style>
            * {
                margin: 0;
                padding: 0;
            }

            html {
                height: 100%;
            }
            
            body {
                height: 100%;
            }

            main {
                display: flex;
                flex-direction: column;
                height: 100%;
            }

            nav {
                background: #000;
                height: 50px;
                padding: 10px;
                display: flex;
                align-items: center;
                border-top: 5px solid #ffd600;
            }

            .logo {
                height: 50px;
            }

            h1 {
                color: #ffd600;
            }

            .content {
                background: #fff;
                flex: 1;
                flex-direction: column;
                height: 100%;
                display: flex;
                justify-content: center;
                align-items: center;
            }

            p {
                font-family: 'Lato', sans-serif;
                color: #333;
                padding: 10px;
                font-size: 20px;
            }
        </style>
        <link href="https://fonts.googleapis.com/css?family=Lato:400" rel="stylesheet" type="text/css">
    </head>
    <body>
        <main>
            <nav>
                <a href="https://npandl.com">
					<img class="logo" src="data:image/png;base64," />
                </a>
            </nav>
            <div class="content">
                <p>Password reset was successful.</p>
                <p>Please check your email for a temporary password to access your account.</p>
            </div>
        </main>
    </body>
</html>
`
