package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/lib/token"
	"github.com/kmhebb/serverExample/log"
)

type Handler struct {
	e             EndpointFunc
	dec           DecodeFunc
	enc           EncodeFunc
	errorFunc     ErrorFunc
	errEncodeFunc ErrorEncodeFunc
}

type HandlerOpts struct {
	// Decoder is the decode function to be used by the handler.
	Decoder DecodeFunc

	// Endpoint is the endpoint function to be used by the handler.
	Endpoint EndpointFunc

	// Encoder is the encode function to be used by the handler. If nil,
	// EncodeJSON will be used.
	Encoder EncodeFunc

	// OnError is the error function to be used by the handler. If nil, LogError
	// will be used.
	OnError ErrorFunc

	// ErrEncoder is the encoder used by the handler to return the error to the client.
	// If nil, the default json error encoder will be used.
	ErrEncoder ErrorEncodeFunc
}

func NewHandler(opts HandlerOpts) *Handler {
	if opts.Encoder == nil {
		opts.Encoder = EncodeJSON
	}
	if opts.OnError == nil {
		opts.OnError = LogError
	}
	if opts.ErrEncoder == nil {
		opts.ErrEncoder = EncodeError
	}

	return &Handler{
		dec:           opts.Decoder,
		e:             opts.Endpoint,
		enc:           opts.Encoder,
		errorFunc:     opts.OnError,
		errEncodeFunc: opts.ErrEncoder,
	}
}

const DefaultResponse = `error:`

// func NewHandlerOld(
// 	endpoint EndpointFunc,
// 	decoder DecodeFunc,
// 	encoder EncodeFunc,
// ) *Handler {
// 	return &Handler{
// 		e:             endpoint,
// 		dec:           decoder,
// 		enc:           encoder,
// 		errorFunc:     DefaultErrorFunc,
// 		errEncodeFunc: DefaultErrorEncodeFunc,
// 	}
// }

// func NewHTMLHandlerDefunct(
// 	endpoint EndpointFunc,
// 	decoder DecodeFunc,
// 	encoder EncodeFunc,
// 	errEncodeFunc ErrorEncodeFunc,
// ) *Handler {
// 	return &Handler{
// 		e:             endpoint,
// 		dec:           decoder,
// 		enc:           encoder,
// 		errorFunc:     DefaultErrorFunc,
// 		errEncodeFunc: errEncodeFunc,
// 	}
// }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := cloud.NewContext(r)
	l := log.NewLogger()

	request, err := h.dec(&ctx, r)
	if err != nil {
		h.errorFunc(ctx, err)
		h.errEncodeFunc(ctx, w, DefaultResponse, err)
		return
	}
	l.Debug(ctx.Ctx, "decoded request", log.Fields{"req": request})

	err = h.authenticate(&ctx)
	if err != nil {
		h.errorFunc(ctx, err)
		h.errEncodeFunc(ctx, w, DefaultResponse, err)
		return
	}

	if ctx.ConfTokenReqired {
		err = h.checkConfirmation(ctx)
		if err != nil {
			h.errorFunc(ctx, err)
			h.errEncodeFunc(ctx, w, DefaultResponse, err)
			return
		}
	}

	response, err := h.e(ctx, request)
	if err != nil {
		h.errorFunc(ctx, err)
		h.errEncodeFunc(ctx, w, response, err)
		return
	}
	l.Debug(ctx.Ctx, "got response", log.Fields{"resp": response})

	if err := h.enc(ctx, w, response); err != nil {
		h.errorFunc(ctx, err)
		h.errEncodeFunc(ctx, w, response, err)
		return
	}
	l.Debug(ctx.Ctx, "encoded response", nil)
}

func (h *Handler) Use(mw EndpointMiddleware) {
	h.e = mw(h.e)
}

// type codedError interface {
// 	Code() int
// }

// type messagedError interface {
// 	Message() string
// }

// type kindedError interface {
// 	Kind() errors.Kind
// }

func (h *Handler) authenticate(ctx *cloud.Context) *cloud.Error {
	l := log.NewLogger()
	hdr := ctx.Token

	if hdr == "" && !ctx.TokenRequired {
		return nil
	}

	if hdr == "" && ctx.TokenRequired {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "token is required",
			//Cause: err,
		}) //fmt.Errorf("bearer token missing or invalid")
	}

	if !strings.HasPrefix(hdr, "Bearer Authorization:") && ctx.TokenRequired {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "bearer token is missing or invalid",
			//Cause: err,
		}) //fmt.Errorf("bearer token missing or invalid")
	}

	t := strings.TrimPrefix(hdr, "Bearer Authorization:")
	if t == "" && ctx.TokenRequired {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "bearer token is missing or invalid",
			//Cause: err,
		}) //fmt.Errorf("bearer token missing or invalid")
	}
	l.Info(ctx.Ctx, "Bearer token present", log.Fields{"token": t})

	claims, err := token.Parse(t)
	if err != nil {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "token error",
			Cause:   err,
		}) //fmt.Errorf("token error: %s", err)
	}
	l.Info(ctx.Ctx, "Parsed bearer token", log.Fields{"user_id": claims.UID})

	// Evaluation of the claims will take place in the service logic as it pertains to the specific request.
	ctx.Token = t
	ctx.UserKey = claims.UID

	return nil

}

func (h *Handler) checkConfirmation(ctx cloud.Context) *cloud.Error {
	if ctx.ConfirmationToken == "" {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "confirmation token is missing or invalid",
			//Cause: err,
		}) // fmt.Errorf("confirmation token missing or invalid")
	}
	// Evaluation of the token will happen in the service.
	// Here we are just making sure there is a confirmation token.
	return nil
}

// Encodes an error response in our standard format. This is the one piece of
// our Handler pipeline that can not be adjusted.
func EncodeError(ctx cloud.Context, w http.ResponseWriter, data interface{}, e *cloud.Error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code())
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"kind":    e.Kind(),
			"message": e.Message(),
		},
	})
}

func EncodeErrorHTML(ctx cloud.Context, w http.ResponseWriter, data interface{}, e *cloud.Error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(e.Code())
	response := fmt.Sprintf("%v\n  kind - %v,\n message: %v", data, e.Kind(), e.Message())
	fmt.Fprintf(w, response)
}
