package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/internal/service"
	"github.com/kmhebb/serverExample/web"
)

type UserService interface {
	// Procedures implemented in the external API, also implemented in the service.
	CreateNewUser(ctx cloud.Context, req service.CreateNewUserRequest) (*service.NewUserResponse, *cloud.Error)
	FindOrCreate(ctx cloud.Context, email, first, last string) (*cloud.User, *cloud.Error)
	Get(ctx cloud.Context, req service.GetUserRequest) (*service.GetUserResponse, *cloud.Error)
	Login(ctx cloud.Context, req service.LoginUserRequest) (*service.LoginUserResponse, *cloud.Error)
	Put(ctx cloud.Context, req service.PutUserRequest) *cloud.Error
	RequestPasswordReset(ctx cloud.Context, req service.UserRequest) *cloud.Error
	ResetPassword(ctx cloud.Context, req service.ResetPasswordRequest) *cloud.Error
	ListUsers(ctx cloud.Context) ([]*cloud.User, *cloud.Error)

	// Internal procedure implmented in the service.
	UpdateLastActivity(ctx cloud.Context, req service.GetUserRequest) *cloud.Error
	ValidateUserAuth(ctx cloud.Context) *cloud.Error

	// Procedures created but unimplemented in the external API.
	ConfirmValidation(ctx cloud.Context, req service.ConfirmValidationRequest) (interface{}, *cloud.Error)
	ValidateEmail(ctx cloud.Context, req service.UserRequest) (interface{}, *cloud.Error)

	// Unimplemented procedures of the service, held for user management development
	DisableUser(ctx cloud.Context, req service.UserRequest) *cloud.Error
	EnableUser(ctx cloud.Context, req service.UserRequest) *cloud.Error
	Find(ctx cloud.Context, req service.UserRequest) (*service.GetUserResponse, *cloud.Error)
}

func RegisterUserRoutes(srv *web.Server, svc service.UserService) {

	routes := map[string]web.HandlerOpts{
		"/users/create": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.CreateNewUserRequest
				//ctx.TokenRequired = true
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode create user request",
						Cause:   err,
					}) //fmt.Errorf("decode create user request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.CreateNewUserRequest)
				return svc.CreateNewUser(ctx, req)
			},
		},
		"/users/get": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.GetUserRequest
				ctx.TokenRequired = true
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode get user by id request",
						Cause:   err,
					}) //fmt.Errorf("decode get user request: %w", err)
				}
				return request, nil

			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.GetUserRequest)
				return svc.Get(ctx, req)
			},
		},
		"/users/login": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.LoginUserRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					fmt.Println(r.Body)
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode login request",
						Cause:   err,
					}) //fmt.Errorf("decode user login request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.LoginUserRequest)
				return svc.Login(ctx, req)
			},
		},
		"/users/put": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.PutUserRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode update user request",
						Cause:   err,
					}) //fmt.Errorf("decode user put request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.PutUserRequest)
				return svc.Put(ctx, req)
			},
		},
		"/users/requestpasswordreset": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.UserRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode get user by email request",
						Cause:   err,
					}) //fmt.Errorf("decode user request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UserRequest)
				return svc.RequestPasswordReset(ctx, req)
			},
		},
		"/users/resetpassword": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.ResetPasswordRequest
				ctx.ConfTokenReqired = true
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode user password reset request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.ResetPasswordRequest)
				return svc.ResetPassword(ctx, req)
			},
			Encoder:    web.EncodeHTML,
			ErrEncoder: web.EncodeErrorHTML,
		},
		"/users/listusers": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				var request service.ListUsersRequest
				//ctx.TokenRequired = true
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode list user request",
						Cause:   err,
					}) //fmt.Errorf("decode list users request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.ListUsersRequest)
				return svc.ListUsers(ctx, req)
			},
		},
		// "/users/find": {
		// 	Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
		// 		var request service.UserRequest
		// 		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// 			return nil, cloud.NewError(cloud.ErrOpts{
		// 				Kind:    cloud.ErrKindBadRequest,
		// 				Message: "failed to decode get user by email request",
		// 				Cause:   err,
		// 			}) //fmt.Errorf("decode user request: %w", err)
		// 		}
		// 		return request, nil
		// 	},
		// 	Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
		// 		req := request.(service.UserRequest)
		// 		return svc.Find(ctx, req)
		// 	},
		// },
	}

	for path, opts := range routes {
		h := web.NewHandler(opts)
		h.Use(web.LoggingMiddleware)
		srv.Handle(path, h)
	}
}
