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

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "CreateUser=Login,Logout,VerificationEmail,VerificationUsername,GetUser,DeleteUser,UpdatePassword,UpdateUsername,UpdateEmail,GetUsers,CreateRecoveryCode,RecoveryPassword"

func errCreateUser(logger *zap.Logger, err error, code int) middleware.Responder {
	if code < http.StatusInternalServerError {
		logger.With(zap.String(log.Error, "client"), zap.Int(log.HTTPStatus, code)).Info(err.Error())
	} else {
		logger.With(zap.String(log.Error, "server"), zap.Int(log.HTTPStatus, code)).Warn(err.Error())
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return operations.NewCreateUserDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}
