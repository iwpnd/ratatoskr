package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/iwpnd/valhalla-builder/services"
	"github.com/iwpnd/valhalla-builder/states"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	blobStore := services.NewBlobStore("blobstore", logger)
	fileStore := services.NewFileStore("filestore", logger)

	name := "baguette"
	args := states.Args{Name: name, Logger: logger, BlobStore: blobStore, FileStore: fileStore}

	// deadline := 1000 * time.Millisecond
	ctx := context.Background()
	// ctx, cancelCtx := context.WithTimeout(ctx, deadline)
	// defer cancelCtx()

	err := states.Execute(ctx, args)
	if err != nil {
		slog.Error("something went wrong", "error", err)
	}
}
