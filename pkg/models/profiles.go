package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/database"
)

type Profile struct {
	UserID    uuid.UUID `json:"userId" validate:"required" minLength:"36" maxLength:"36" example:"00000000-0000-0000-0000-000000000000"`
	Username  string    `json:"username" validate:"required" minLength:"2" maxLength:"64" example:"Krusher99"`
	AvatarURI *string   `json:"avatarUri" validate:"optional" maxLength:"1024" example:"http://cdn.discord.com/some-image-id"`
	CreatedAt int64     `json:"createdAt" validate:"required" example:"1700000000"`
	UpdatedAt int64     `json:"updatedAt" validate:"required" example:"1700000000"`
}

type NewProfile struct {
	Username  string  `json:"username" validate:"required" minLength:"2" maxLength:"64" example:"Krusher99"`
	AvatarURI *string `json:"avatarUri,omitempty" validate:"optional" maxLength:"1024" example:"http://cdn.discord.com/some-image-id"`
}

func UpsertProfile(ctx context.Context, q *database.Queries, userID uuid.UUID, profile NewProfile) (Profile, error) {
	var upserted Profile

	row, err := q.UpsertProfile(ctx, database.UpsertProfileParams{
		UserUuid: userID,
		Username: profile.Username,
		Avatar:   profile.AvatarURI,
	})
	if err != nil {
		return upserted, convertSqlErr(err)
	}

	upserted.Username = row.Username
	upserted.AvatarURI = row.Avatar
	upserted.CreatedAt = row.CreatedAt
	upserted.UpdatedAt = row.UpdatedAt
	return upserted, nil
}

func GetProfile(ctx context.Context, q *database.Queries, userID uuid.UUID) (Profile, error) {
	var profile Profile

	row, err := q.GetProfile(ctx, userID)
	if err != nil {
		return profile, convertSqlErr(err)
	}

	profile.Username = row.Username
	profile.AvatarURI = row.Avatar
	profile.CreatedAt = row.CreatedAt
	profile.UpdatedAt = row.UpdatedAt
	return profile, nil
}
