package tiles

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Builder interface {
	BuildConfig(context.Context) error
	BuildTiles(context.Context) error
	BuildTilesExtract(context.Context) error
	BuildAdmins(context.Context) error
	Path() string
	AdminPath() string
	ExtractPath() string
	TilesPath() string
}

type TileBuilder struct {
	concurrency  int
	maxCacheSize int64

	path        string
	tilesPath   string
	extractPath string
	configPath  string
	adminPath   string
	datasetPath string

	configCreated bool

	logger   *slog.Logger
	executor *executor
}

type TileBuilderOptions struct {
	Debug bool

	MaxCacheSize int64
	Concurrency  int
	Path         string
	Dataset      string
}

func createPathIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func requiredBinaries() []string {
	return []string{
		"valhalla_build_admins",
		"valhalla_build_config",
		"valhalla_build_extract",
		"valhalla_build_tiles",
	}
}

func NewTileBuilder(
	opts *TileBuilderOptions,
	logger *slog.Logger) (*TileBuilder, error) {

	executor := &executor{
		logger: logger,
		debug:  opts.Debug,
	}

	for _, b := range requiredBinaries() {
		if !executor.hasExecutable(b) {
			return nil, fmt.Errorf("missing executable: %s", b)
		}
	}

	builder := &TileBuilder{executor: executor, logger: logger}
	builder.concurrency = 2
	if opts.Concurrency != 0 {
		builder.concurrency = opts.Concurrency
	}

	builder.maxCacheSize = 700 * 1048576 // 700MiB
	if opts.MaxCacheSize != 0 {
		builder.maxCacheSize = opts.MaxCacheSize
	}

	if opts.Path != "" {
		if opts.Dataset == "" {
			return nil, fmt.Errorf("set path but missing dataset name")
		}

		err := createPathIfNotExists(opts.Path)
		if err != nil {
			return nil, fmt.Errorf("error creating basepath: %d", err)
		}
		builder.path = opts.Path
		builder.tilesPath = builder.path + "/valhalla_tiles"
		err = createPathIfNotExists(builder.tilesPath)
		if err != nil {
			return nil, fmt.Errorf("error creating valhalla_tiles path: %s", err)
		}
		builder.extractPath = builder.path + "/valhalla_tiles.tar"
		builder.adminPath = builder.path + "/admin.sqlite"
		builder.configPath = builder.path + "/config.json"
		builder.datasetPath = builder.path + "/" + opts.Dataset
		return builder, nil
	}

	return builder, nil
}

func (ve *TileBuilder) BuildConfig(ctx context.Context) error {
	params := []string{
		"--mjolnir-concurrency", fmt.Sprint(ve.concurrency),
		"--mjolnir-max-cache-size", fmt.Sprint(ve.maxCacheSize),
		"--mjolnir-tile-dir", ve.tilesPath,
		"--mjolnir-tile-extract", ve.extractPath,
		"--mjolnir-admin", ve.adminPath,
	}

	ve.logger.Info("creating valhalla config", "params", params)

	output, err := ve.executor.executeWithOutput(ctx, "valhalla_build_config", params)
	if err != nil {
		return err
	}

	err = os.WriteFile(ve.configPath, output, 0644)
	if err != nil {
		return fmt.Errorf("error creating valhalla config: %s", err)
	}

	if _, err := os.Stat(ve.configPath); os.IsNotExist(err) {
		return fmt.Errorf("error creating valhalla config: %s", err)
	}

	ve.configCreated = true

	ve.logger.Info(
		"finished creating valhalla config",
		"params", params,
	)

	return nil
}

func (ve *TileBuilder) BuildTiles(ctx context.Context) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	params := []string{"--config", ve.configPath, ve.datasetPath}
	ve.logger.Info("started creating tiles", "params", params)

	err := ve.executor.execute(ctx, "valhalla_build_tiles", params)
	if err != nil {
		return err
	}

	ve.logger.Info(
		"finished creating tiles",
		"params", params,
	)

	return nil
}

func (ve *TileBuilder) BuildTilesExtract(ctx context.Context) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	start := time.Now()

	params := []string{"--config", ve.configPath, "-O"}
	ve.logger.Info("started tarballing tiles extract", "params", params)

	err := ve.executor.execute(ctx, "valhalla_build_extract", params)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	ve.logger.Info(
		"finished tarballing tile extract",
		"params", params,
		"elapsed", elapsed.String(),
	)

	return nil
}

func (ve *TileBuilder) BuildAdmins(ctx context.Context) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	params := []string{"--config", ve.configPath, ve.datasetPath}
	ve.logger.Info("started building admins", "params", params)

	err := ve.executor.execute(ctx, "valhalla_build_admins", params)
	if err != nil {
		return err
	}

	ve.logger.Info(
		"finished building admins",
		"params", params,
	)

	return nil
}

func (ve *TileBuilder) Path() string {
	return ve.path
}

func (ve *TileBuilder) ExtractPath() string {
	return ve.extractPath
}

func (ve *TileBuilder) AdminPath() string {
	return ve.adminPath
}

func (ve *TileBuilder) TilesPath() string {
	return ve.tilesPath
}
