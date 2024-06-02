package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/iwpnd/valhalla-builder/services/download"
	"github.com/iwpnd/valhalla-builder/services/tiles"
	"github.com/iwpnd/valhalla-builder/states"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	name := "europe/andorra"

	path := "./bin/builder/data"

	builderOpts := &tiles.TileBuilderOptions{
		Path:         path,
		Concurrency:  6,
		MaxCacheSize: 700 * 1048576, // 700MiB
		Debug:        false,
		Dataset:      "andorra-latest.osm.pbf",
	}
	builder, err := tiles.NewTileBuilder(builderOpts, logger)
	if err != nil {
		logger.Error("cannot instantiate executor: ", "err", err)
		panic(err)
	}

	downloaderOpts := download.GeofabrikDownloaderOptions{
		Dataset:    name,
		OutputPath: path,
	}
	downloader := download.NewGeofabrikDownloader(downloaderOpts)

	args := states.Args{
		Name:       name,
		Logger:     logger,
		Builder:    builder,
		Downloader: downloader,
	}

	err = states.Execute(ctx, args)
	if err != nil {
		slog.Error("something went wrong", "error", err)
	}
}
