package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/providers"
	"github.com/opendungeon/opendungeon/internal/services"
	"github.com/opendungeon/opendungeon/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, disableUserCreation bool, db *services.DB, email string, password string) (uuid.UUID, error) {
	if disableUserCreation {
		return uuid.Nil, fiber.NewError(fiber.StatusForbidden, "User creation is disabled.")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "Password could not be encrypted.")
	}
	passwordDigest := string(bytes)

	user, err := models.CreateUser(ctx, db.Queries, models.NewUser{
		Email: email,
	})
	if err != nil {
		if errors.Is(err, models.ErrCheckViolation) {
			return uuid.Nil, fiber.ErrBadRequest
		}
		if errors.Is(err, models.ErrUniqueViolation) {
			return uuid.Nil, fiber.ErrConflict
		}

		log.Errorf("failed to create user: %v", err)
		return uuid.Nil, fiber.ErrInternalServerError
	}

	_, err = models.CreateIdentity(ctx, db.Queries, user.ID, models.NewIdentity{
		Provider:       "email",
		PasswordDigest: &passwordDigest,
	})
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create identity.")
	}

	return user.ID, nil
}

func SignIn(ctx context.Context, db *services.DB, email string, password string) (uuid.UUID, error) {
	identity, err := models.GetIdentity(ctx, db.Queries, email, "email")
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusNotFound, "Failed to find identity.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*identity.PasswordDigest), []byte(password)); err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusNotFound, "Failed to find identity.")
	}

	return identity.UserID, nil
}

type CallbackRedirect struct {
	UserID   uuid.UUID
	Redirect *url.URL
}

func DiscordCallback(
	ctx context.Context,
	disableUserCreation bool,
	db *services.DB,
	clientID, clientSecret string,
	baseUrl, clientUrl *url.URL,
	code, state string,
) (CallbackRedirect, error) {
	var cr CallbackRedirect

	discord := providers.NewDiscord(baseUrl, clientID, clientSecret)

	token, err := discord.Exchange(ctx, code)
	if err != nil {
		log.Errorf("failed to exchange auth code with discord: %v", err)
		return cr, fiber.NewError(fiber.StatusPreconditionFailed, "Failed to sign in with discord.")
	}

	discordUser, err := discord.GetUserInfo(ctx, token)
	if err != nil {
		log.Errorf("failed to get user info from discord: %v", err)
		return cr, fiber.NewError(fiber.StatusPreconditionFailed, "Failed to get account info from Discord.")
	}

	// HANDLE EXISTING DISCORD IDENTITY
	identity, err := models.GetIdentity(ctx, db.Queries, discordUser.Email, "discord")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("failed to retrieve identity: %v", err)
		return cr, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve identity.")
	}

	identityExists := identity.ProviderUid != nil && *identity.ProviderUid == discordUser.ID
	if identityExists {
		cr.UserID = identity.UserID
		cr.Redirect = clientUrl // redirect to home page '/'
		return cr, nil
	}

	// HANDLE EXISTING USER
	existingUser, err := db.Queries.GetUserByEmail(ctx, discordUser.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("failed to retrieve existing user: %v", err)
		return cr, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve user.")
	}

	userExists := existingUser.Email == discordUser.Email
	if userExists {
		_, err = models.CreateIdentity(ctx, db.Queries, existingUser.Uuid, models.NewIdentity{
			Provider:    "discord",
			ProviderUid: &discordUser.ID,
		})
		if err != nil {
			log.Errorf("failed to create identity on existing user: %v", err)
			return cr, fiber.NewError(fiber.StatusInternalServerError, "Failed to create identity.")
		}

		cr.UserID = existingUser.Uuid
		cr.Redirect = clientUrl // redirect to home page '/'
		return cr, nil
	}

	// HANDLE CREATING A NEW USER
	if disableUserCreation {
		return cr, fiber.NewError(fiber.StatusForbidden, "User creation is disabled.")
	}

	user, err := models.CreateUser(ctx, db.Queries, models.NewUser{Email: discordUser.Email})
	if err != nil {
		// no reason to check for database errors here, since the email MUST be unique as
		// we already checked if it exists, AND it must be valid since it came from discord
		log.Errorf("failed to create new user during discord sign in: %v", err)
		return cr, fiber.NewError(fiber.StatusInternalServerError, "Failed to create user.")
	}

	_, err = models.CreateIdentity(ctx, db.Queries, existingUser.Uuid, models.NewIdentity{
		Provider:    "discord",
		ProviderUid: &discordUser.ID,
	})
	if err != nil {
		log.Errorf("failed to create new identity during discord sign in: %v", err)
		return cr, fiber.NewError(fiber.StatusInternalServerError, "Failed to create identity.")
	}

	_, err = models.UpsertProfile(ctx, db.Queries, user.ID, models.NewProfile{
		Username:  discordUser.Username,
		AvatarURI: discordUser.AvatarUri,
	})
	if err != nil {
		log.Warn("failed to create profile for discord user: %v", err)
	}

	cr.UserID = user.ID
	cr.Redirect = clientUrl.JoinPath("/profiles/me/edit")
	return cr, nil
}
