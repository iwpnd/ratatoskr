package download

import (
	"context"
)

// Downloader interface
type Downloader interface {
	Get(ctx context.Context, dataset string, outputPath string) error
	MD5(ctx context.Context, dataset string) (string, error)
}
