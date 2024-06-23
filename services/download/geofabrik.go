package download

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/iwpnd/go-geofabrik"
)

type GeofabrikDownloaderOptions struct {
	Url        string
	Dataset    string
	OutputPath string
}

type GeofabrikDownloader struct {
	opts GeofabrikDownloaderOptions
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, fs.ErrNotExist)
}

func NewGeofabrikDownloader(
	opts GeofabrikDownloaderOptions,
) *GeofabrikDownloader {
	return &GeofabrikDownloader{opts: opts}
}

func (od *GeofabrikDownloader) Get(ctx context.Context) error {
	g, err := geofabrik.New(od.opts.Url)
	if err != nil {
		return fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	err = g.Download(ctx, od.opts.Dataset, od.opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed download from geofabrik dataset %s: %w", od.opts.Dataset, err)
	}

	return nil
}
