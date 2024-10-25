package states

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iwpnd/valhalla-builder/services/compress"
	"github.com/iwpnd/valhalla-builder/services/download"
	"github.com/iwpnd/valhalla-builder/services/tiles"
)

type Params struct {
	Name string

	Downloader download.Downloader
	Builder    tiles.Builder
	Compressor compress.Compressor

	Logger *slog.Logger
}

func (p Params) validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if p.Name == "" {
		return fmt.Errorf("during Args validation: name cannot be empty string")
	}
	if p.Logger == nil {
		return fmt.Errorf("during Args validation: logger cannot be nil")
	}
	if p.Builder == nil {
		return fmt.Errorf("during Args validation: builder cannot be nil")
	}

	return nil
}
