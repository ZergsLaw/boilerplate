package rest

import (
	"errors"
	"net"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/boilerplate/internal/app"
)

func (svc *service) verificationEmail(params operations.VerificationEmailParams) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.userApp.VerificationEmail(ctx, string(params.Args.Email))
	switch {
	case err == nil:
		return operations.NewVerificationEmailNoContent()
	case errors.Is(err, app.ErrEmailExist):
		return errVerificationEmail(log, err, http.StatusConflict)
	default:
		return errVerificationEmail(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) verificationUsername(params operations.VerificationUsernameParams) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.userApp.VerificationUsername(ctx, string(params.Args.Username))
	switch {
	case err == nil:
		return operations.NewVerificationUsernameNoContent()
	case errors.Is(err, app.ErrUsernameExist):
		return errVerificationUsername(log, err, http.StatusConflict)
	default:
		return errVerificationUsername(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) createUser(params operations.CreateUserParams) middleware.Responder {
	ctx, log, remoteIP := fromRequest(params.HTTPRequest, nil)

	origin := app.Origin{
		IP:        net.ParseIP(remoteIP),
		UserAgent: params.HTTPRequest.Header.Get("User-Agent"),
	}

	u, token, err := svc.userApp.CreateUser(
		ctx,
		string(params.Args.Email),
		string(params.Args.Username),
		string(params.Args.Password),
		origin,
	)
	switch {
	case err == nil:
		cookie := generateCookie(token)
		return operations.NewCreateUserOK().WithPayload(User(u)).WithSetCookie(cookie.String())
	case errors.Is(err, app.ErrEmailExist):
		return errCreateUser(log, err, http.StatusConflict)
	case errors.Is(err, app.ErrUsernameExist):
		return errCreateUser(log, err, http.StatusConflict)
	default:
		return errCreateUser(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) login(params operations.LoginParams) middleware.Responder {
	ctx, log, remoteIP := fromRequest(params.HTTPRequest, nil)

	origin := app.Origin{
		IP:        net.ParseIP(remoteIP),
		UserAgent: params.HTTPRequest.Header.Get("User-Agent"),
	}

	u, token, err := svc.userApp.Login(ctx, string(params.Args.Email), string(params.Args.Password), origin)
	switch {
	case err == nil:
		cookie := generateCookie(token)
		return operations.NewLoginOK().WithPayload(User(u)).WithSetCookie(cookie.String())
	case errors.Is(err, app.ErrNotFound):
		return errLogin(log, err, http.StatusNotFound)
	case errors.Is(err, app.ErrNotValidPassword):
		return errLogin(log, err, http.StatusBadRequest)
	default:
		return errLogin(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) logout(params operations.LogoutParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.userApp.Logout(ctx, *authUser)
	switch {
	case err == nil:
		return operations.NewLogoutNoContent()
	default:
		return errLogout(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) getUser(params operations.GetUserParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	getUserID := authUser.ID
	if params.ID != nil {
		getUserID = app.UserID(*params.ID)
	}

	u, err := svc.userApp.User(ctx, *authUser, getUserID)
	switch {
	case err == nil:
		return operations.NewGetUserOK().WithPayload(User(u))
	case errors.Is(err, app.ErrNotFound):
		return errGetUser(log, err, http.StatusNotFound)
	default:
		return errGetUser(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) deleteUser(params operations.DeleteUserParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.userApp.DeleteUser(ctx, *authUser)
	switch {
	case err == nil:
		return operations.NewDeleteUserNoContent()
	default:
		return errDeleteUser(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) updatePassword(params operations.UpdatePasswordParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.userApp.UpdatePassword(ctx, *authUser, string(params.Args.Old), string(params.Args.New))
	switch {
	case err == nil:
		return operations.NewUpdatePasswordNoContent()
	case errors.Is(err, app.ErrNotValidPassword):
		return errUpdatePassword(log, err, http.StatusConflict)
	default:
		return errUpdatePassword(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) updateUsername(params operations.UpdateUsernameParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.userApp.UpdateUsername(ctx, *authUser, string(params.Args.Username))
	switch {
	case err == nil:
		return operations.NewUpdateUsernameNoContent()
	case errors.Is(err, app.ErrUsernameExist):
		return errUpdateUsername(log, err, http.StatusConflict)
	case errors.Is(err, app.ErrUsernameNeedDifferentiate):
		return errUpdateUsername(log, err, http.StatusConflict)
	default:
		return errUpdateUsername(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) updateEmail(params operations.UpdateEmailParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.userApp.UpdateEmail(ctx, *authUser, string(params.Args.Email))
	switch {
	case err == nil:
		return operations.NewUpdateEmailNoContent()
	case errors.Is(err, app.ErrEmailExist):
		return errUpdateEmail(log, err, http.StatusConflict)
	case errors.Is(err, app.ErrEmailNeedDifferentiate):
		return errUpdateEmail(log, err, http.StatusConflict)
	default:
		return errUpdateEmail(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) createRecoveryCode(params operations.CreateRecoveryCodeParams) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.userApp.CreateRecoveryCode(ctx, string(params.Args.Email))
	switch {
	case err == nil:
		return operations.NewCreateRecoveryCodeNoContent()
	case errors.Is(err, app.ErrNotFound):
		return errCreateRecoveryCode(log, err, http.StatusNotFound)
	default:
		return errCreateRecoveryCode(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) recoveryPassword(params operations.RecoveryPasswordParams) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.userApp.RecoveryPassword(ctx, string(params.Args.RecoveryCode), string(params.Args.Password))
	switch {
	case err == nil:
		return operations.NewRecoveryPasswordNoContent()
	case errors.Is(err, app.ErrNotFound):
		return errRecoveryPassword(log, err, http.StatusNotFound)
	case errors.Is(err, app.ErrCodeExpired):
		return errRecoveryPassword(log, err, http.StatusBadRequest)
	default:
		return errRecoveryPassword(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) getUsers(params operations.GetUsersParams, authUser *app.AuthUser) middleware.Responder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	page := app.Page{
		Limit:  int(params.Limit),
		Offset: int(swag.Int32Value(params.Offset)),
	}

	u, total, err := svc.userApp.ListUserByUsername(ctx, *authUser, params.Username, page)
	switch {
	case err == nil:
		return operations.NewGetUsersOK().WithPayload(&operations.GetUsersOKBody{
			Total: swag.Int32(int32(total)),
			Users: Users(u),
		})
	default:
		return errGetUsers(log, err, http.StatusInternalServerError)
	}
}
