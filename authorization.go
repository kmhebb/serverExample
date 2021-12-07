package cloud

const DefaultTokenLength = 32

type Action string

const (
// ActionAccessTideAPI Action = "tideapi.access"
)

type AuthorizationService interface {
	Authorize(ctx Context, act Action) error
}

type AuthorizationFunc func(ctx Context, act Action) error
