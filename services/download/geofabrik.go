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
	Dataset    string
	OutputPath string
}

type GeofabrikDownloader struct {
	baseUrl string
	opts    GeofabrikDownloaderOptions
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, fs.ErrNotExist)
}

func NewGeofabrikDownloader(
	opts GeofabrikDownloaderOptions,
) *GeofabrikDownloader {
	return &GeofabrikDownloader{opts: opts, baseUrl: "http://download.geofabrik.de"}
}

func (od *GeofabrikDownloader) Get(ctx context.Context) error {
	g, err := geofabrik.New(od.baseUrl)
	if err != nil {
		return fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	err = g.Download(ctx, od.opts.Dataset, od.opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed download from geofabrik dataset %s: %w", od.opts.Dataset, err)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (od *GeofabrikDownloader) MD5(ctx context.Context) (string, error) {
	g, err := geofabrik.New(od.baseUrl)
	if err != nil {
		return "", fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	md5, err := g.MD5(ctx, od.opts.Dataset)
	if err != nil {
		return "", fmt.Errorf("failed fetch MD5 from geofabrik dataset %s: %w", od.opts.Dataset, err)
	}

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return md5, nil

}
