package download

import (
	"context"
)

type Downloader interface {
	Get(ctx context.Context) error
	MD5(ctx context.Context) (string, error)
}
