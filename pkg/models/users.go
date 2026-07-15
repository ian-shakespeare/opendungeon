package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/database"
)

type User struct {
	ID      uuid.UUID `json:"id" validate:"required" minLength:"36" maxLength:"36" example:"00000000-0000-0000-0000-000000000000"`
	Email   string    `json:"email" validate:"required" minLength:"5" maxLength:"255" example:"john.doe@email.com"`
	IsAdmin bool      `json:"isAdmin" validate:"required" default:"false"`
}

type NewUser struct {
	Email string `json:"email" validate:"required" minLength:"5" maxLength:"255" example:"john.doe@email.com"`
}

func CreateUser(ctx context.Context, q *database.Queries, user NewUser) (User, error) {
	var created User

	id := uuid.New()
	row, err := q.CreateUser(ctx, database.CreateUserParams{
		Uuid:  id,
		Email: user.Email,
	})
	if err != nil {
		return created, convertSqlErr(err)
	}

	created.ID = row.Uuid
	created.Email = row.Email
	created.IsAdmin = row.IsAdmin
	return created, nil
}

func GetUser(ctx context.Context, q *database.Queries, email string) (User, error) {
	var user User

	row, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		return user, convertSqlErr(err)
	}

	user.ID = row.Uuid
	user.Email = row.Email
	user.IsAdmin = row.IsAdmin
	return user, nil
}
