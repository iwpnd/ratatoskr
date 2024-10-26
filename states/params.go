package states

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iwpnd/valhalla-tiles-builder/services/compress"
	"github.com/iwpnd/valhalla-tiles-builder/services/download"
	"github.com/iwpnd/valhalla-tiles-builder/services/tiles"
)

type Params struct {
	Name string

	Downloader download.Downloader
	Builder    tiles.Builder
	Compressor compress.Compressor

	Logger *slog.Logger
}

func (p Params) Validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if p.Name == "" {
		return fmt.Errorf("during Params validation: name cannot be empty string")
	}
	if p.Logger == nil {
		return fmt.Errorf("during Params validation: logger cannot be nil")
	}

	return nil
}
