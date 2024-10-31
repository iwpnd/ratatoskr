package compress

import "context"

// Compressor interface
type Compressor interface {
	Compress(ctx context.Context, archive string, files ...string) error
}
