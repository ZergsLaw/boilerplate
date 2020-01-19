package rest

import (
	"github.com/go-openapi/swag"
	"github.com/sirupsen/logrus"
	"github.com/zergslaw/users/internal/api/rest/generated/models"
	"github.com/zergslaw/users/internal/api/rest/generated/restapi/operations"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "CreateUser=Login,Logout,VerificationEmail,VerificationUsername,GetUser,DeleteUser,UpdatePassword,UpdateUsername,UpdateEmail,GetUsers"

//nolint:dupl,goconst
func errCreateUser(log logrus.FieldLogger, err error, code int) operations.CreateUserResponder { //nolint:deadcode,unused
	if code < 500 {
		log.WithFields(logrus.Fields{LogHTTPStatus: code, LogError: "client"}).Info(err)
	} else {
		log.WithFields(logrus.Fields{LogHTTPStatus: code, LogError: "server"}).Warn(err)
	}

	msg := err.Error()
	if code == 500 { // Do no expose details about internal errors.
		msg = "internal error"
	}

	return operations.NewCreateUserDefault(code).WithPayload(&models.Error{Message: swag.String(msg)})
}
