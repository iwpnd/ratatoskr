package download

import (
	"context"
)

type Downloader interface {
	Get(ctx context.Context) error
}
