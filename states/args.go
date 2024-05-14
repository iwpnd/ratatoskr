package states

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iwpnd/valhalla-builder/services"
)

type Args struct {
	Name string

	BlobStore services.BlobStorer
	FileStore services.FileStorer

	Logger *slog.Logger
}

func (a Args) validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if a.Name == "" {
		return fmt.Errorf("name cannot be empty string")
	}
	if a.BlobStore == nil {
		return fmt.Errorf("services cannot be nil")
	}
	if a.FileStore == nil {
		return fmt.Errorf("services cannot be nil")
	}
	if a.Logger == nil {
		return fmt.Errorf("services cannot be nil")
	}

	return nil
}
