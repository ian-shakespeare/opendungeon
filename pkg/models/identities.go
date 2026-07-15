package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/opendungeon/opendungeon/internal/database"
)

type Identity struct {
	UserID         uuid.UUID `json:"userId" validate:"required" minLength:"36" maxLength:"36" example:"00000000-0000-0000-0000-000000000000"`
	Provider       string    `json:"provider" validate:"required" enums:"[email,discord]" example:"discord"`
	ProviderUid    *string   `json:"providerUid" validate:"optional" maxLength:"255"`
	PasswordDigest *string   `json:"passwordDigest" validate:"optional" minLength:"60" maxLength:"discord"`
}

type NewIdentity struct {
	Provider       string  `json:"provider" validate:"required" enums:"[email,discord]" example:"email"`
	ProviderUid    *string `json:"providerUid" validate:"optional" maxLength:"255"`
	PasswordDigest *string `json:"passwordDigest" validate:"optional" minLength:"60" maxLength:"60"`
}

func CreateIdentity(ctx context.Context, q *database.Queries, userID uuid.UUID, identity NewIdentity) (Identity, error) {
	var created Identity

	row, err := q.CreateIdentity(ctx, database.CreateIdentityParams{
		UserUuid:       userID,
		Provider:       identity.Provider,
		ProviderUid:    identity.ProviderUid,
		PasswordDigest: identity.PasswordDigest,
	})
	if err != nil {
		return created, convertSqlErr(err)
	}

	created.UserID = row.Uuid
	created.Provider = row.Name
	created.ProviderUid = row.ProviderUid
	created.PasswordDigest = row.PasswordDigest
	return created, nil
}

func ListIdentities(ctx context.Context, q *database.Queries, email string) ([]Identity, error) {
	rows, err := q.ListIdentities(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Identity{}, nil
		}

		return nil, convertSqlErr(err)
	}

	identities := make([]Identity, 0, len(rows))
	for _, row := range rows {
		identities = append(identities, Identity{
			UserID:         row.UserUuid,
			Provider:       row.Provider,
			ProviderUid:    row.ProviderUid,
			PasswordDigest: row.PasswordDigest,
		})
	}

	return identities, nil
}

func GetIdentity(ctx context.Context, q *database.Queries, email, provider string) (Identity, error) {
	var identity Identity

	row, err := q.GetIdentityByEmail(ctx, database.GetIdentityByEmailParams{
		Email:    email,
		Provider: provider,
	})
	if err != nil {
		return identity, convertSqlErr(err)
	}

	identity.UserID = row.UserUuid
	identity.Provider = row.Provider
	identity.ProviderUid = row.ProviderUid
	identity.PasswordDigest = row.PasswordDigest
	return identity, nil
}
