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
	Url           string
	Dataset       string
	OutputPath    string
	ForceDownload bool
}

type GeofabrikDownloader struct {
	opts GeofabrikDownloaderOptions
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

	filePath := od.opts.OutputPath + "/" + od.opts.Dataset
	if _, err := os.Stat(filePath); errors.Is(err, fs.ErrNotExist) {
		if !od.opts.ForceDownload {
			return nil
		}
	}

	err = g.Download(ctx, od.opts.Dataset, od.opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed download from geofabrik dataset %s: %w", od.opts.Dataset, err)
	}

	return nil
}
