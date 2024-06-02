package stores

import (
	"context"
	"fmt"
	"log/slog"
)

type BlobStorer interface {
	Get(ctx context.Context, id string) error
}

type BlobStore struct {
	name   string
	logger *slog.Logger
}

func NewBlobStore(name string, logger *slog.Logger) *FileStore {
	return &FileStore{name: name, logger: logger}
}

func (b *BlobStore) Get(ctx context.Context, id string) error {
	msg := fmt.Sprintf("getting something with %s", b.name)

	b.logger.Info(
		msg,
		"id", id,
		"service", b.name,
	)

	return nil
}
