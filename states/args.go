package states

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iwpnd/valhalla-builder/services/download"
	"github.com/iwpnd/valhalla-builder/services/tiles"
)

type Args struct {
	Name string

	Downloader *download.GeofabrikDownloader
	Builder    *tiles.TileBuilder

	Logger *slog.Logger
}

func (a Args) validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if a.Name == "" {
		return fmt.Errorf("name cannot be empty string")
	}
	if a.Downloader == nil {
		return fmt.Errorf("services cannot be nil")
	}
	if a.Builder == nil {
		return fmt.Errorf("services cannot be nil")
	}
	if a.Logger == nil {
		return fmt.Errorf("services cannot be nil")
	}

	return nil
}
