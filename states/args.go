package states

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iwpnd/valhalla-builder/services/compress"
	"github.com/iwpnd/valhalla-builder/services/download"
	"github.com/iwpnd/valhalla-builder/services/tiles"
)

type Args struct {
	Name string

	Downloader download.Downloader
	Builder    tiles.Builder
	Compressor compress.Compressor

	Logger *slog.Logger
}

func (a Args) validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if a.Name == "" {
		return fmt.Errorf("during Args validation: name cannot be empty string")
	}
	if a.Downloader == nil {
		return fmt.Errorf("during Args validation: services cannot be nil")
	}
	if a.Builder == nil {
		return fmt.Errorf("during Args validation: services cannot be nil")
	}
	if a.Logger == nil {
		return fmt.Errorf("during Args validation: services cannot be nil")
	}

	return nil
}
