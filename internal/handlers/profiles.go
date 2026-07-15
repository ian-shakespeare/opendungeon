package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/services"
	"github.com/opendungeon/opendungeon/pkg/models"
)

func UpsertProfile(ctx context.Context, db *services.DB, userID uuid.UUID, profile models.NewProfile) (models.Profile, error) {
	upserted, err := models.UpsertProfile(ctx, db.Queries, userID, profile)
	if err != nil {
		if errors.Is(err, models.ErrCheckViolation) {
			return upserted, fiber.ErrBadRequest
		}
		if errors.Is(err, models.ErrNotFound) {
			return upserted, fiber.ErrNotFound
		}

		log.Errorf("failed to upsert profile: %v", err)
		return upserted, fiber.NewError(fiber.StatusInternalServerError, "Failed to create profile.")
	}

	return upserted, err
}

func GetProfile(ctx context.Context, db *services.DB, userID uuid.UUID) (models.Profile, error) {
	profile, err := models.GetProfile(ctx, db.Queries, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return profile, fiber.ErrNotFound
		}

		log.Errorf("failed to get profile: %v", err)
		return profile, fiber.ErrInternalServerError
	}

	return profile, nil
}
