package compress

import "context"

type Compressor interface {
	Compress(ctx context.Context, archive string, files ...string) error
}
