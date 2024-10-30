package states

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/iwpnd/ratatoskr/services/compress"
	"github.com/iwpnd/ratatoskr/services/download"
	"github.com/iwpnd/ratatoskr/services/tiles"
)

type Params struct {
	dataset    string
	basePath   string
	outputPath string

	logger *slog.Logger

	downloader download.Downloader
	builder    tiles.Builder
	compressor compress.Compressor
}

func NewParams(dataset string, basePath string, logger *slog.Logger) *Params {
	return &Params{
		dataset:  dataset,
		basePath: basePath,
		logger:   logger,
	}
}

func (p *Params) setOutputPath(dataset string, md5 string) error {
	path := []string{
		strings.TrimRightFunc(
			p.basePath, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			},
		),
	}

	if dataset != "" {
		path = append(path, dataset)
	}

	if md5 != "" {
		path = append(path, md5)
	}

	p.outputPath = strings.Join(path, "/")

	return nil
}

func (p *Params) Validate(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if p.dataset == "" {
		return fmt.Errorf("during Params validation: dataset cannot be empty string")
	}
	if p.basePath == "" {
		return fmt.Errorf("during Params validation: basePath cannot be empty string")
	}
	if p.logger == nil {
		return fmt.Errorf("during Params validation: logger cannot be nil")
	}
	if p.builder == nil {
		return fmt.Errorf("during Params validation: tiles builder cannot be nil")
	}
	if p.downloader == nil {
		return fmt.Errorf("during Params validation: downloader cannot be nil")
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
