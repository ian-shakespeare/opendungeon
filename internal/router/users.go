package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/opendungeon/opendungeon/internal/handlers"
	"github.com/opendungeon/opendungeon/pkg/models"
)

// createUser
//
//	@Summary		Create user
//	@Description	Create a new user.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.NewUser	true	"New user information"
//	@Success		201		{object}	models.User
//	@Failure		409		{string}	string	"Conflict"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/api/users [post]
func (r *router) createUser(c fiber.Ctx) error {
	var user models.NewUser
	if err := c.Bind().JSON(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	created, err := handlers.CreateUser(c.Context(), r.db, r.disableUserCreation, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}
