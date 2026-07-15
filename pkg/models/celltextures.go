package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/opendungeon/opendungeon/internal/database"
)

type CellTexture struct {
	Key         string `json:"key" validate:"required" minLength:"2" maxLength:"64" example:"tundra.grass"`
	DisplayName string `json:"displayName" validate:"required" minLength:"2" maxLength:"64" example:"Tundra Grass"`
	CreatedAt   int64  `json:"createdAt" validate:"required" example:"1700000000"`
	UpdatedAt   int64  `json:"updatedAt" validate:"required" example:"1700000000"`
}

type NewCellTexture struct {
	Key         string `json:"key" validate:"required" minLength:"2" maxLength:"64" example:"tundra.grass"`
	DisplayName string `json:"displayName" validate:"required" minLength:"2" maxLength:"64" example:"Tundra Grass"`
}

func CreateCellTexture(ctx context.Context, q *database.Queries, cellTexture NewCellTexture) (CellTexture, error) {
	var created CellTexture

	row, err := q.CreateCellTexture(ctx, database.CreateCellTextureParams{
		Key:         cellTexture.Key,
		DisplayName: cellTexture.DisplayName,
	})
	if err != nil {
		return created, convertSqlErr(err)
	}

	created.Key = row.Key
	created.DisplayName = row.DisplayName
	created.CreatedAt = row.CreatedAt
	created.UpdatedAt = row.UpdatedAt
	return created, nil
}

func ListCellTextures(ctx context.Context, q *database.Queries) ([]CellTexture, error) {
	rows, err := q.ListCellTextures(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []CellTexture{}, nil
		}

		return nil, convertSqlErr(err)
	}

	ct := make([]CellTexture, 0, len(rows))
	for _, row := range rows {
		ct = append(ct, CellTexture{
			Key:         row.Key,
			DisplayName: row.DisplayName,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return ct, nil
}
