package download

import (
	"context"
	"fmt"

	"github.com/iwpnd/go-geofabrik"
)

type GeofabrikDownloaderOptions struct {
	Dataset    string
	OutputPath string
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
	g, err := geofabrik.New("http://download.geofabrik.de", false)
	if err != nil {
		return fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	err = g.Download(ctx, od.opts.Dataset, od.opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed download from geofabrik dataset %s: %w", od.opts.Dataset, err)

	}

	return nil
}
