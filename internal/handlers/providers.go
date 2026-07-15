package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/opendungeon/opendungeon/internal/providers"
	"github.com/opendungeon/opendungeon/pkg/models"
)

func ListProviders(ctx context.Context, redirectUri *url.URL, discordClientID, discordClientSecret string) (models.ProviderList, error) {
	var pl models.ProviderList

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Errorf("failed to generate state: %v", err)
		return pl, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate state.")
	}
	pl.State = hex.EncodeToString(b)

	if discordClientID != "" && discordClientSecret != "" {
		discord := providers.NewDiscord(redirectUri, discordClientID, discordClientSecret)
		pl.Providers = append(pl.Providers, models.Provider{
			Name:    "Discord",
			AuthURL: discord.AuthUrl(pl.State),
		})
	}

	return pl, nil
}
