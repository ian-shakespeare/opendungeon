package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/opendungeon/opendungeon/internal/handlers"
)

// listProviders
//
//	@Summary		List providers
//	@Description	List available auth providers.
//	@Tags			Providers
//	@Produce		json
//	@Success		200	{object}	[]models.Provider	"Available auth providers"
//	@Failure		500	{string}	string				"Server error"
//	@Router			/api/providers [get]
func (r *router) listProviders(c fiber.Ctx) error {
	redirectUri := r.clientURL.JoinPath("discord", "callback")
	providers, err := handlers.ListProviders(c.Context(), redirectUri, r.discordClientID, r.discordClientSecret)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    providers.State,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.JSON(providers.Providers)
}
