package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/log"
)

// DecodeFunc is a function for decoding information from an HTTP request into a
// concrete type used by business logic.
type DecodeFunc func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error)

// EncodeFunc is a function for encoding business logic results into a response
// to an HTTP request.
type EncodeFunc func(ctx cloud.Context, w http.ResponseWriter, response interface{}) *cloud.Error

// EndpointFunc is a function for performing business logic using concrete
// request and response types that are respectively provided by a DecodeFunc
// and used by an EncodeFund.
type EndpointFunc func(ctx cloud.Context, request interface{}) (response interface{}, err *cloud.Error)

// ErrorFunc is a function called when an error is returned by any of the other
// function types defined in this package: DecodeFunc, EncodeFunc, or
// EndpointFunc. It's responsible internal handling of the error only, the error
// response to the associated HTTP request is automatic.
type ErrorFunc func(cloud.Context, *cloud.Error)

// ErrorEncodeFunc encodes an error for a particular client. JSON is standard, html should be inserted as necessary.
type ErrorEncodeFunc func(ctx cloud.Context, w http.ResponseWriter, data interface{}, e *cloud.Error)

// ContentType[x] is the value to use in the encode functions.
const ContentTypeJSON = "application/json; charset=utf-8"
const ContentTypeXML = "application/xml; charset=utf-8"
const ContentTypeHTML = "text/html; charset=utf-8"

// EncodeJSON is an EncodeFunc that responds with a 200 OK status code and a
// body that contains the response value marshaled to JSON.
func EncodeJSON(ctx cloud.Context, w http.ResponseWriter, response interface{}) *cloud.Error {
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return cloud.NewError(cloud.ErrOpts{Cause: err})
	}
	return nil
}

func EncodeXML(ctx cloud.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", ContentTypeXML)
	w.WriteHeader(http.StatusOK)
	return xml.NewEncoder(w).Encode(response)
}

func EncodeHTML(ctx cloud.Context, w http.ResponseWriter, data interface{}) *cloud.Error {
	w.Header().Set("Content-Type", ContentTypeHTML)
	response := fmt.Sprint(data)
	fmt.Fprintf(w, response)
	return nil
}

// LogError is an ErrorFunc that simply logs the error.
func LogError(ctx cloud.Context, e *cloud.Error) {
	l := log.NewLogger()
	l.Error(ctx.Ctx, e, string(e.Kind()), nil)
}

// var DefaultErrorFunc = func(ctx cloud.Context, err error) {
// 	l := log.NewLogger()
// 	msg := "unexpected error"
// 	if kinder, ok := err.(kindedError); ok {
// 		msg = string(kinder.Kind())
// 	}
// 	l.Error(ctx.Ctx, err, msg, nil)
// }

// var DefaultErrorEncodeFunc = func(ctx cloud.Context, w http.ResponseWriter, response interface{}, err error) error {
// 	cause := goerrs.Unwrap(err)

// 	code := http.StatusInternalServerError
// 	msg := errors.DefaultMessage
// 	kind := errors.EINTERNAL

// 	if coder, ok := cause.(codedError); ok {
// 		code = coder.Code()
// 	}
// 	if messager, ok := cause.(messagedError); ok {
// 		msg = messager.Message()
// 	}
// 	if kinder, ok := cause.(kindedError); ok {
// 		kind = kinder.Kind()
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"error": map[string]interface{}{
// 			"kind":    kind,
// 			"message": msg,
// 		},
// 	})
// 	return nil
// }

// func EncodeWebError(ctx cloud.Context, err error) {
// 	l := log.NewLogger()
// 	msg := "unexpected error"
// 	if kinder, ok := err.(kindedError); ok {
// 		msg = string(kinder.Kind())
// 	}
// 	l.Error(ctx.Ctx, err, msg, nil)
// }
