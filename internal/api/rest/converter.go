package rest

import (
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/models"
	"github.com/zergslaw/boilerplate/internal/app"
)

// Users conversion []app.User => []*models.User.
func Users(u []app.User) []*models.User {
	users := make([]*models.User, len(u))

	for i := range users {
		users[i] = User(&u[i])
	}

	return users
}

// User conversion app.User => models.User.
func User(u *app.User) *models.User {
	return &models.User{
		ID:       models.UserID(u.ID),
		Username: models.Username(u.Username),
		Email:    models.Email(u.Email),
	}
}
