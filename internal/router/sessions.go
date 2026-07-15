package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/opendungeon/opendungeon/internal/handlers"
	"github.com/opendungeon/opendungeon/pkg/models"
)

// createEmailSession
//
//	@Summary		Create session with email
//	@Description	Create a new session using email authentication.
//	@Tags			Sessions
//	@Accept			json
//	@Param			request	body	models.NewEmailSession	true	"Email credentials"
//	@Success		201		"Session id cookie"
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/api/sessions/email [post]
func (r *router) createEmailSession(c fiber.Ctx) error {
	var s models.NewEmailSession

	if err := c.Bind().JSON(&s); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := handlers.CreateSessionWithEmail(c.Context(), r.db, s.Email, s.Password)
	if err != nil {
		return err
	}

	sess := session.FromContext(c)
	sess.Set("user_id", user.ID)
	return c.SendStatus(fiber.StatusCreated)
}

// createDiscordSession
//
//	@Summary		Create session with Discord
//	@Description	Create a new session using Discord authentication.
//	@Description	If a user with the associated email does not exist, it is created.
//	@Tags			Sessions
//	@Accept			json
//	@Param			request	body	models.NewDiscordSession	true	"Discord credentials"
//	@Success		201		"Session id cookie"
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/api/sessions/discord [post]
func (r *router) createDiscordSession(c fiber.Ctx) error {
	var s models.NewDiscordSession

	if err := c.Bind().JSON(&s); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	redirectUri := r.clientURL.JoinPath("discord", "callback")
	user, err := handlers.CreateSessionWithDiscord(c.Context(), r.db, r.discordClientID, r.discordClientSecret, redirectUri, s.AuthCode)
	if err != nil {
		return err
	}

	sess := session.FromContext(c)
	sess.Set("user_id", user.ID)
	return c.SendStatus(fiber.StatusCreated)
}

// deleteSession
//
//	@Summary		Delete session
//	@Description	Delete the active session.
//	@Tags			Sessions
//	@Success		204
//	@Failure		401	{string}	string	"Unauthorized"
//	@Router			/api/sessions [delete]
func (r *router) deleteSession(c fiber.Ctx) error {
	sess := session.FromContext(c)
	sess.Delete("user_id")
	return c.SendStatus(fiber.StatusNoContent)
}
