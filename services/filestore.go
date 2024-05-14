package services

import (
	"context"
	"fmt"
	"log/slog"
)

type FileStorer interface {
	Get(ctx context.Context, id string) error
}

type FileStore struct {
	name   string
	logger *slog.Logger
}

func NewFileStore(name string, logger *slog.Logger) *FileStore {
	return &FileStore{name: name, logger: logger}
}

func (f *FileStore) Get(ctx context.Context, id string) error {
	msg := fmt.Sprintf("getting something with %s", f.name)

	f.logger.Info(
		msg,
		"id", id,
		"service", f.name,
	)

	return nil
}
