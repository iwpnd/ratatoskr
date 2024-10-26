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
	name   string
	logger *slog.Logger

	downloader download.Downloader
	builder    tiles.Builder
	compressor compress.Compressor
}

func NewParams(name string, logger *slog.Logger) *Params {
	return &Params{
		name:   name,
		logger: logger,
	}
}

func (p *Params) Validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if p.name == "" {
		return fmt.Errorf("during Params validation: name cannot be empty string")
	}
	if p.logger == nil {
		return fmt.Errorf("during Params validation: logger cannot be nil")
	}
	if p.compressor != nil && p.builder == nil {
		return fmt.Errorf("cannot use compressor without tiles builder")
	}

	return nil
}

func (p *Params) WithDownload(downloader download.Downloader) *Params {
	p.downloader = downloader
	return p
}

func (p *Params) WithTileBuilder(builder tiles.Builder) *Params {
	p.builder = builder
	return p
}

func (p *Params) WithCompression(compressor compress.Compressor) *Params {
	p.compressor = compressor
	return p
}
