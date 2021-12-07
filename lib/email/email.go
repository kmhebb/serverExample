// Package email provides a no-op implementation of EmailService for use in testing.
package email

import (
	cloud "github.com/kmhebb/serverExample"
)

type Service interface {
	//NewCustomerAsync(ctx cloud.Context)
	NewUserAsync(ctx cloud.Context, name, to, pass string)
	ResetPasswordAsync(ctx cloud.Context, name, to, token string)
	NewPasswordAsync(ctx cloud.Context, name, to, pass string)
	ValidateEmailAsync(ctx cloud.Context, to, code string)
	Close() chan int
	TestConnection() error
}

func NewNoOpService() Service {
	return noOpService{}
}

type noOpService struct{}

func (s noOpService) NewCustomerAsync(ctx cloud.Context)                                 {}
func (s noOpService) NewPasswordAsync(ctx cloud.Context, name, email, passphrase string) {}
func (s noOpService) NewUserAsync(ctx cloud.Context, name, email, passphrase string)     {}
func (s noOpService) ResetPasswordAsync(ctx cloud.Context, name, email, token string)    {}
func (s noOpService) ValidateEmailAsync(ctx cloud.Context, to, code string)              {}
func (s noOpService) TestConnection() error {
	return nil
}

// BUG: Should the done channel be closed? Do we care in a testing-only environment?
func (s noOpService) Close() chan int {
	// We create a buffered channel and then fill it so that the client gets a
	// value immediately.
	done := make(chan int, 1)
	done <- 1
	return done
}
