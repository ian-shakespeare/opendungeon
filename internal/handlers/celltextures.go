package handlers

import (
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/opendungeon/opendungeon/internal/services"
	"github.com/opendungeon/opendungeon/pkg/models"
)

const (
	CellTextureWidth  = 64
	CellTextureHeight = 64
)

func CreateCellTexture(
	ctx context.Context,
	db *services.DB,
	storage *services.Storage,
	key, displayName string,
	content io.Reader,
) (models.CellTexture, error) {
	var created models.CellTexture

	if len(key) < 3 || 64 < len(key) {
		return created, fiber.NewError(http.StatusBadRequest, "Key must be between 3 and 64 (inclusive) characters in length.")
	}

	if len(displayName) < 3 || 64 < len(displayName) {
		return created, fiber.NewError(http.StatusBadRequest, "Display name must be between 3 and 64 (inclusive) characters in length.")
	}

	img, err := png.Decode(content)
	if err != nil {
		return created, fiber.NewError(http.StatusUnsupportedMediaType, "Image must be a PNG format.")
	}

	rect := img.Bounds()
	width := rect.Max.X
	height := rect.Max.Y
	if width != CellTextureWidth || height != CellTextureHeight {
		message := fmt.Sprintf("Image must have a width of %d pixels and a height of %d pixels.", CellTextureWidth, CellTextureHeight)
		return created, fiber.NewError(http.StatusBadRequest, message)
	}

	created, err = models.CreateCellTexture(ctx, db.Queries, models.NewCellTexture{
		Key:         key,
		DisplayName: displayName,
	})
	if err != nil {
		if errors.Is(err, models.ErrUniqueViolation) {
			return created, fiber.ErrConflict
		}

		log.Errorf("failed to create cell texture record: %v", err)
		return created, fiber.NewError(http.StatusInternalServerError, "Failed to create texture record.")
	}

	// use a pipe to avoid creating another buffer
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		_ = png.Encode(pw, img) // know this wont error since we decoded from a PNG
	}()

	scopedKey := "celltexture." + created.Key
	if _, err := storage.CreateFile(scopedKey, "image/png", pr); err != nil {
		// clean up db entry since the actual file didn't make it. ignore errors since we can't do anything about it.
		_, _ = db.Queries.HardDeleteCellTexture(ctx, scopedKey)

		log.Errorf("failed to store cell texture: %v", err)
		return created, fiber.NewError(http.StatusInternalServerError, "Failed to store file.")
	}

	return created, nil
}

func GetCellTexture(
	ctx context.Context,
	db *services.DB,
	storage *services.Storage,
	key string,
) (io.ReadCloser, error) {
	texture, err := db.Queries.GetCellTexture(ctx, key)
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "Texture not found.")
	}

	scopedKey := "celltexture." + texture.Key
	reader, err := storage.GetFile(scopedKey)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "Failed to retrieve file.")
	}

	return reader, nil
}

func ListCellTextures(
	ctx context.Context,
	db *services.DB,
) ([]models.CellTexture, error) {
	textures, err := models.ListCellTextures(ctx, db.Queries)
	if err != nil {
		log.Errorf("failed to list textures: %v", err)
		return nil, fiber.NewError(http.StatusInternalServerError, "Failed to list textures.")
	}

	return textures, nil
}
