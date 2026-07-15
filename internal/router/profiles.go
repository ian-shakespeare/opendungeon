package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/opendungeon/opendungeon/internal/handlers"
	"github.com/opendungeon/opendungeon/pkg/models"
)

// upsertMyProfile
//
//	@Summary		Create or replace user's profile
//	@Description	Create or replace the profile for the authenticated user.
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			request		body	models.NewProfile	true	"Profile data"
//	@Success		201		{object}	models.Profile
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/api/profiles/me [put]
func (r *router) upsertMyProfile(c fiber.Ctx) error {
	userId, ok := getUserId(c)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var profile models.NewProfile
	if err := c.Bind().JSON(&profile); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body.")
	}

	upserted, err := handlers.UpsertProfile(c.Context(), r.db, userId, profile)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(upserted)
}

// getMyProfile
//
//	@Summary		Get user's profile
//	@Description	Get the profile for the authenticated user.
//	@Tags			Profiles
//	@Produce		json
//	@Success		200	{object}	models.Profile
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/api/profiles/me [get]
func (r router) getMyProfile(c fiber.Ctx) error {
	userId, ok := getUserId(c)
	if !ok {
		return fiber.ErrUnauthorized
	}

	profile, err := handlers.GetProfile(c.Context(), r.db, userId)
	if err != nil {
		return err
	}

	return c.JSON(profile)
}
