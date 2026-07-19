package handlers

import (
	"context"
	"errors"
	"io"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/database"
	"github.com/opendungeon/opendungeon/internal/media"
	"github.com/opendungeon/opendungeon/internal/services"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UpsertedProfile struct {
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

func UpsertProfile(
	ctx context.Context,
	db *services.DB,
	storage *services.Storage,
	userId uuid.UUID,
	username string,
	avatar io.Reader,
) (database.UpsertProfileRow, error) {
	var avatarID *string
	if avatar != nil {
		converted, err := media.ConvertToAvatar(avatar)
		if err != nil {
			if errors.Is(err, media.ErrUnknownContentType) || errors.Is(err, media.ErrUnsupportedImageFormat) {
				return database.UpsertProfileRow{}, fiber.NewError(fiber.StatusBadRequest, "Invalid avatar format. Must be a PNG, JPEG, HEIC, or WEBP.")
			}

			log.Errorf("failed to convert avatar: %v", err)
			return database.UpsertProfileRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to convert avatar.")
		}

		id := uuid.New()
		scopedKey := "avatar." + id.String()
		if _, err := storage.CreateFile(scopedKey, "image/png", converted); err != nil {
			return database.UpsertProfileRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to save avatar.")
		}

		idStr := id.String()
		avatarID = &idStr
	}

	upserted, err := db.Queries.UpsertProfile(ctx, database.UpsertProfileParams{
		UserUuid: userId,
		Username: username,
		Avatar:   avatarID,
	})
	if err != nil {
		if avatarID != nil {
			scopedKey := "avatar." + *avatarID
			_ = storage.DeleteFile(scopedKey)
		}

		sqlErr := new(sqlite.Error)
		if errors.As(err, &sqlErr) {
			if sqlErr.Code() == sqlite3.SQLITE_CONSTRAINT_CHECK {
				return database.UpsertProfileRow{}, fiber.NewError(fiber.StatusBadRequest, "Invalid request.")
			}
		}
		return database.UpsertProfileRow{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to create profile.")
	}

	return upserted, err
}

func GetProfile(ctx context.Context, db *services.DB, userId uuid.UUID) (database.GetProfileRow, error) {
	profile, err := db.Queries.GetProfile(ctx, userId)
	if err != nil {
		return profile, fiber.NewError(fiber.StatusNotFound, "Profile not found.")
	}

	return profile, nil
}
