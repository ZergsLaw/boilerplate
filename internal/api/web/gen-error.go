// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package web

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/zergslaw/boilerplate/internal/api/web/generated/models"
	"github.com/zergslaw/boilerplate/internal/api/web/generated/restapi/operations"
	"github.com/zergslaw/boilerplate/internal/log"
	"go.uber.org/zap"
)

func errLogin(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewLoginDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errLogout(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewLogoutDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errVerificationEmail(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewVerificationEmailDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errVerificationUsername(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewVerificationUsernameDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errGetUser(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewGetUserDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errDeleteUser(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewDeleteUserDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errUpdatePassword(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewUpdatePasswordDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errUpdateUsername(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewUpdateUsernameDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errUpdateEmail(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewUpdateEmailDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errGetUsers(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewGetUsersDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errCreateRecoveryCode(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewCreateRecoveryCodeDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}

func errRecoveryPassword(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewRecoveryPasswordDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}