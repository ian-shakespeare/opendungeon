package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/opendungeon/opendungeon/internal/services"
	"github.com/opendungeon/opendungeon/pkg/models"
)

func CreateUser(ctx context.Context, db *services.DB, disableUserCreation bool, user models.NewUser) (models.User, error) {
	if disableUserCreation {
		return models.User{}, fiber.ErrBadRequest
	}

	created, err := models.CreateUser(ctx, db.Queries, user)
	if err != nil {
		if errors.Is(err, models.ErrUniqueViolation) {
			return created, fiber.ErrConflict
		}

		log.Errorf("failed to create user: %v", err)
		return created, fiber.ErrInternalServerError
	}

	return created, nil
}
