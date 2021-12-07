package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pborman/uuid"
)

// type Context struct {
// 	Ctx       context.Context
// 	Logger    log.Logger
// 	Request   *http.Request
// 	RequestID uuid.UUID
// }

type Context struct {
	// Ctx is equivalent to Request.Context() and is provided as a convenience.
	// If you are attaching cancellations or timers to the context, you should
	// always make a copy of the this context instead of attempting to modify it
	// directly, as those changes will not propagate backwards through the call
	// chain.
	Ctx context.Context

	// Request is the original HTTP request.
	Request *http.Request

	// RequestID is a UUID set on every incoming request. This may be set
	// directly by the client or it may be set by the server if the client did
	// not provide one. The request ID provides a single identifier that is used
	// throughout our call stack and our instrumentation to correlate  log
	// entries, error reports, APM, etc.
	RequestID string

	// UserKey is the user id of the requesting user.
	UserKey string

	// Many requests will be accompanied by a token. We will include this in the context to make it easy to access.
	Token string

	// Many endpoints require a token. This variable will be set in the decode func so that the auth middleware can know.
	TokenRequired bool

	// Some endpoints require a confirmation key or token. This is different than the API token and needs a different path for evaluation.
	ConfirmationToken string

	// Is the ConfirmationToken required?
	ConfTokenReqired bool
}

// NewContext returns a Context with fields parsed from a source HTTP request.
func NewContext(r *http.Request) Context {
	fmt.Printf("new context, token: %v", r.Header.Get("Authorization"))
	ctx := Context{
		Ctx:               r.Context(),
		Request:           r,
		RequestID:         uuid.New(),
		Token:             r.Header.Get("Authorization"),
		TokenRequired:     false,
		ConfirmationToken: r.URL.Query().Get("conf"),
		ConfTokenReqired:  false,
	}

	return ctx
}
