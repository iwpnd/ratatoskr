# Ratatoskr

The only thing that is somewhat of a hurdle to building your own
routing microservices using the
[valhalla routing engine](https://github.com/valhalla/valhalla) is the
creation and maintenance of valhalla routing tiles.

Ratatoskr is a pipeline to download osm data, build admins, build valhalla
routing tiles, compress the result and upload to blob storage. It can be used
as a one off, or tied to a task workers, cronjobs or what have you.

## Requirement

Requires you to have the valhalla executables in the environment you
are running ratatoskr in.

```bash
valhalla_build_config
valhalla_build_tiles
valhalla_build_extract
valhalla_build_admin
```

## Installation

### cli

```bash
NAME:
   ratatoskr run - run valhalla tiles pipeline

USAGE:
   ratatoskr run [command [command options]]

DESCRIPTION:
   Run the entire valhalla tiles pipeline

           1. Download the dataset
           2. Build a valhalla configuration file
           3. Build valhalla routing tiles
           4. Build valhalla routing tiles .tar file
           5. Compress valhalla routing tiles .tar file
           6. (optional, wip) upload to Blob Storage


OPTIONS:
   --outputpath value, -o value
   --dataset value, -d value
   --build.concurrency value     (default: 0)
   --build.maxCacheSize value    (default: 0)
   --help, -h                    show help
```

### package

```bash
go get -u github.com/iwpnd/ratatoskr
```

```go
package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/iwpnd/ratatoskr/pipeline"
    "github.com/iwpnd/ratatoskr/services/compress"
    "github.com/iwpnd/ratatoskr/services/download"
    "github.com/iwpnd/ratatoskr/services/tiles"
    "github.com/iwpnd/ratatoskr/states"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()

    dataset := "europe/germany/bremen"
    path := "./tmp"

    options := []tiles.Option{
        tiles.WithConcurrency(runtime.NumCPU()),
        tiles.WithMaxCacheSize(1024 * 1048576),
    }

    builder, err := tiles.NewTileBuilder(logger, options...)
    if err != nil {
        logger.Error("cannot instantiate executor: ", "err", err)
        os.Exit(1)
    }

    downloader := download.NewGeofabrikDownloader()
    compressor := &compress.GzipCompressor{}

    params := states.NewParams(dataset, path, logger).
        WithDownload(downloader).
        WithTileBuilder(builder).
        WithCompression(compressor)

    err = pipeline.Execute(ctx, params)
    if err != nil {
        logger.Error("something went wrong", "error", err)
        os.Exit(1)
    }
}
```

## License

MIT

## Acknowledgement

## Maintainer

Benjamin Ramser - [@iwpnd](https://github.com/iwpnd)

Project Link: [https://github.com/iwpnd/ratatoskr](https://github.com/iwpnd/ratatoskr)
