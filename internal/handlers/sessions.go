package handlers

import (
	"context"
	"errors"
	"net/url"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/opendungeon/opendungeon/internal/providers"
	"github.com/opendungeon/opendungeon/internal/services"
	"github.com/opendungeon/opendungeon/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateSessionWithEmail(ctx context.Context, db *services.DB, email string, password string) (models.User, error) {
	user, err := models.GetUser(ctx, db.Queries, email)
	if err != nil {
		return user, fiber.ErrNotFound
	}

	identities, err := models.ListIdentities(ctx, db.Queries, email)
	if err != nil {
		log.Errorf("failed to list identities: %v", err)
		return user, fiber.ErrInternalServerError
	}

	// no identities. create a new one.
	if len(identities) == 0 {
		b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Errorf("failed to hash password: %v", err)
			return user, fiber.ErrInternalServerError
		}
		digest := string(b)

		_, err = models.CreateIdentity(ctx, db.Queries, user.ID, models.NewIdentity{
			Provider:       "email",
			PasswordDigest: &digest,
		})
		if err != nil {
			log.Errorf("failed to create identity: %v", err)
			return user, fiber.ErrInternalServerError
		}

		return user, nil
	}

	identityIndex := slices.IndexFunc(identities, func(i models.Identity) bool {
		return i.Provider == "email"
	})
	if identityIndex == -1 {
		return user, fiber.ErrNotFound
	}

	identity := identities[identityIndex]
	if err := bcrypt.CompareHashAndPassword([]byte(*identity.PasswordDigest), []byte(password)); err != nil {
		return user, fiber.ErrNotFound
	}

	return user, nil
}

func CreateSessionWithDiscord(
	ctx context.Context,
	db *services.DB,
	clientID, clientSecret string,
	redirectUri *url.URL,
	code string,
) (models.User, error) {
	// TODO: oauth state
	var user models.User

	if clientID == "" || clientSecret == "" {
		return user, fiber.ErrBadRequest
	}

	discord := providers.NewDiscord(redirectUri, clientID, clientSecret)

	token, err := discord.Exchange(ctx, code)
	if err != nil {
		log.Errorf("failed to exchange: %v", err)
		return user, fiber.ErrPreconditionFailed
	}

	userInfo, err := discord.GetUserInfo(ctx, token)
	if err != nil {
		log.Errorf("failed to get user info: %v", err)
		return user, fiber.ErrPreconditionFailed
	}

	user, err = models.GetUser(ctx, db.Queries, userInfo.Email)
	if errors.Is(err, models.ErrNotFound) {
		user, err = models.CreateUser(ctx, db.Queries, models.NewUser{
			Email: userInfo.Email,
		})
	}
	if err != nil {
		log.Errorf("failed to get user: %v", err)
		return user, fiber.ErrInternalServerError
	}

	identity, err := models.GetIdentity(ctx, db.Queries, user.Email, "discord")
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		log.Errorf("failed to get identity: %v", err)
		return user, fiber.ErrInternalServerError
	}

	isWrongDiscordUser := !errors.Is(err, models.ErrNotFound) && userInfo.ID != *identity.ProviderUid
	if isWrongDiscordUser {
		return user, fiber.ErrNotFound
	}

	_, err = models.CreateIdentity(ctx, db.Queries, user.ID, models.NewIdentity{
		Provider:    "discord",
		ProviderUid: &userInfo.ID,
	})
	if err != nil {
		log.Errorf("failed to create identity: %v", err)
		return user, fiber.ErrInternalServerError
	}

	return user, nil
}
