package download

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/iwpnd/go-geofabrik"
)

type GeofabrikDownloader struct {
	baseUrl string
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, fs.ErrNotExist)
}

func NewGeofabrikDownloader() *GeofabrikDownloader {
	return &GeofabrikDownloader{baseUrl: "http://download.geofabrik.de"}
}

func (od *GeofabrikDownloader) Get(ctx context.Context, dataset string, outputPath string) error {
	g, err := geofabrik.New(od.baseUrl)
	if err != nil {
		return fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	err = g.Download(ctx, dataset, outputPath)
	if err != nil {
		return fmt.Errorf("failed download from geofabrik dataset %s: %w", dataset, err)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (od *GeofabrikDownloader) MD5(ctx context.Context, dataset string) (string, error) {
	g, err := geofabrik.New(od.baseUrl)
	if err != nil {
		return "", fmt.Errorf("cannot instantiate geofabrik: %w", err)
	}

	md5, err := g.MD5(ctx, dataset)
	if err != nil {
		return "", fmt.Errorf("failed fetch MD5 from geofabrik dataset %s: %w", dataset, err)
	}

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return md5, nil

}
