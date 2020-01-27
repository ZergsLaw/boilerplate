package rest

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/sirupsen/logrus"
	"github.com/zergslaw/users/internal/api/rest/generated/models"
	"github.com/zergslaw/users/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/users/internal/log"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "CreateUser=Login,Logout,VerificationEmail,VerificationUsername,GetUser,DeleteUser,UpdatePassword,UpdateUsername,UpdateEmail,GetUsers"

//nolint:dupl,goconst
func errCreateUser(logger logrus.FieldLogger, err error, code int) operations.CreateUserResponder { //nolint:deadcode,unused
	if code < http.StatusInternalServerError {
		logger.WithFields(logrus.Fields{log.HTTPStatus: code, log.Error: "client"}).Info(err)
	} else {
		logger.WithFields(logrus.Fields{log.HTTPStatus: code, log.Error: "server"}).Warn(err)
	}

	msg := err.Error()
	if code == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = "internal error"
	}

	return operations.NewCreateUserDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}
