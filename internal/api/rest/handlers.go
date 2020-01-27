package rest

import (
	"errors"
	"net"
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/zergslaw/users/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/users/internal/app"
)

func (svc *service) verificationEmail(params operations.VerificationEmailParams) operations.VerificationEmailResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.app.VerificationEmail(ctx, string(params.Email))
	switch {
	case err == nil:
		return operations.NewVerificationEmailNoContent()
	case errors.Is(err, app.ErrEmailExist):
		return errVerificationEmail(log, err, http.StatusConflict)
	default:
		return errVerificationEmail(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) verificationUsername(params operations.VerificationUsernameParams) operations.VerificationUsernameResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	err := svc.app.VerificationUsername(ctx, string(params.Username))
	switch {
	case err == nil:
		return operations.NewVerificationUsernameNoContent()
	case errors.Is(err, app.ErrUsernameExist):
		return errVerificationUsername(log, err, http.StatusConflict)
	default:
		return errVerificationUsername(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) createUser(params operations.CreateUserParams) operations.CreateUserResponder {
	ctx, log, remoteIP := fromRequest(params.HTTPRequest, nil)

	origin := app.Origin{
		IP:        net.ParseIP(remoteIP),
		UserAgent: params.HTTPRequest.Header.Get("User-Agent"),
	}

	u, token, err := svc.app.CreateUser(
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

func (svc *service) Login(params operations.LoginParams) operations.LoginResponder {
	ctx, log, remoteIP := fromRequest(params.HTTPRequest, nil)

	origin := app.Origin{
		IP:        net.ParseIP(remoteIP),
		UserAgent: params.HTTPRequest.Header.Get("User-Agent"),
	}

	u, token, err := svc.app.Login(ctx, string(params.Args.Email), string(params.Args.Password), origin)
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

func (svc *service) logout(params operations.LogoutParams, authUser *app.AuthUser) operations.LogoutResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.app.Logout(ctx, *authUser)
	switch {
	case err == nil:
		return operations.NewLogoutNoContent()
	default:
		return errLogout(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) getUser(params operations.GetUserParams, authUser *app.AuthUser) operations.GetUserResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, nil)

	u, err := svc.app.User(ctx, *authUser, app.UserID(params.ID))
	switch {
	case err == nil:
		return operations.NewGetUserOK().WithPayload(User(u))
	case errors.Is(err, app.ErrNotFound):
		return errGetUser(log, err, http.StatusNotFound)
	default:
		return errGetUser(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) deleteUser(params operations.DeleteUserParams, authUser *app.AuthUser) operations.DeleteUserResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.app.DeleteUser(ctx, *authUser)
	switch {
	case err == nil:
		return operations.NewDeleteUserNoContent()
	default:
		return errDeleteUser(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) updatePassword(params operations.UpdatePasswordParams, authUser *app.AuthUser) operations.UpdatePasswordResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.app.UpdatePassword(ctx, *authUser, string(params.Args.Old), string(params.Args.New))
	switch {
	case err == nil:
		return operations.NewUpdatePasswordNoContent()
	case errors.Is(err, app.ErrNotValidPassword):
		return errUpdatePassword(log, err, http.StatusConflict)
	default:
		return errUpdatePassword(log, err, http.StatusInternalServerError)
	}
}

func (svc *service) updateUsername(params operations.UpdateUsernameParams, authUser *app.AuthUser) operations.UpdateUsernameResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.app.UpdateUsername(ctx, *authUser, string(params.Username))
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

func (svc *service) updateEmail(params operations.UpdateEmailParams, authUser *app.AuthUser) operations.UpdateEmailResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	err := svc.app.UpdateEmail(ctx, *authUser, string(params.Email))
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

func (svc *service) getUsers(params operations.GetUsersParams, authUser *app.AuthUser) operations.GetUsersResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, authUser)

	page := app.Page{
		Limit:  int(*params.Args.Pagination.Limit),
		Offset: int(*params.Args.Pagination.Offset),
	}

	u, total, err := svc.app.ListUserByUsername(ctx, *authUser, string(params.Args.Username), page)
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
