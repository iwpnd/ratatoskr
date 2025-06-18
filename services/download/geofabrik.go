package download

import (
	"context"
	"fmt"

	"github.com/iwpnd/go-geofabrik"
)

// GeofabrikDownloader struct.
type GeofabrikDownloader struct {
	baseUrl string
}

// NewGeofabrikDownloader to instantiate a geofarbik downloader.
func NewGeofabrikDownloader() *GeofabrikDownloader {
	return &GeofabrikDownloader{baseUrl: "http://download.geofabrik.de"}
}

// Get to download an OSM dataset to output path.
func (od *GeofabrikDownloader) Get(ctx context.Context, dataset, outputPath string) error {
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

// MD5 to return an OSM dataset MD5 checksum
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
