package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/iwpnd/ratatoskr/pipeline"
	"github.com/iwpnd/ratatoskr/services/compress"
	"github.com/iwpnd/ratatoskr/services/download"
	"github.com/iwpnd/ratatoskr/services/tiles"
	"github.com/iwpnd/ratatoskr/states"
	"github.com/urfave/cli/v3"
)

var runCommand cli.Command

var logger *slog.Logger

func run(ctx context.Context, cmd *cli.Command) error {
	outputPath := cmd.String("outputpath")
	dataset := cmd.String("dataset")

	buildConcurrency := cmd.Int("build.concurrency")
	buildMaxCacheSize := cmd.Int("build.maxCacheSize")

	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	options := []tiles.Option{}
	if buildConcurrency != 0 {
		if buildConcurrency > runtime.NumCPU() {
			buildConcurrency = runtime.NumCPU()

			options = append(options, tiles.WithConcurrency(buildConcurrency))
		}
	}

	if buildMaxCacheSize != 0 {
		options = append(options, tiles.WithMaxCacheSizeInBytes(int64(buildMaxCacheSize)))
	}

	builder, err := tiles.NewTileBuilder(logger, options...)
	if err != nil {
		logger.Error("cannot instantiate executor: ", "err", err)
		return err
	}

	downloader := download.NewGeofabrikDownloader()
	compressor := &compress.GzipCompressor{}

	params := states.NewParams(dataset, outputPath, logger).
		WithDownload(downloader).
		WithTileBuilder(builder).
		WithCompression(compressor)

	err = pipeline.Execute(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	runCommand = cli.Command{
		Name:  "run",
		Usage: "run valhalla tiles pipeline",
		Description: `Run the entire valhalla tiles pipeline

        1. Download the dataset
        2. Build a valhalla configuration file
        3. Build valhalla routing tiles
        4. Build valhalla routing tiles .tar file
        5. Compress valhalla routing tiles .tar file
        6. (optional, wip) upload to Blob Storage
        `,
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "outputpath",
				DefaultText: "test",
				Aliases:     []string{"o"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:     "dataset",
				Aliases:  []string{"d"},
				Required: true,
			},
			&cli.IntFlag{
				Name:     "build.concurrency",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "build.maxCacheSize",
				Required: false,
			},
		},
	}

	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	cmd := &cli.Command{
		Name:  "ratatoskr",
		Usage: "build valhalla routing engine routing tiles",
		Commands: []*cli.Command{
			&runCommand,
		},
	}

	ctx := context.Background()

	if err := cmd.Run(ctx, os.Args); err != nil {
		logger.Error("an error occurred", "error", err)
		os.Exit(1)
	}
}
