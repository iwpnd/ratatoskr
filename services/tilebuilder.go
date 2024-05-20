package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

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

func NewTileBuilder(
	opts *TileBuilderOptions,
	logger *slog.Logger) (*TileBuilder, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("missing path input")
	}

	if opts.Dataset == "" {
		return nil, fmt.Errorf("missing dataset input")
	}

	executor := &executor{
		logger: logger,
		debug:  opts.Debug,
	}
	if !executor.hasExecutable("valhalla_build_tiles") {
		return nil, fmt.Errorf("missing executable: valhalla_build_tiles")
	}

	if !executor.hasExecutable("valhalla_build_extract") {
		return nil, fmt.Errorf("missing executable: valhalla_build_extract")
	}

	if !executor.hasExecutable("valhalla_build_admins") {
		return nil, fmt.Errorf("missing executable: valhalla_build_admins")
	}

	if !executor.hasExecutable("valhalla_build_config") {
		return nil, fmt.Errorf("missing executable: valhalla_build_config")
	}

	builder := &TileBuilder{executor: executor, logger: logger}
	err := createPathIfNotExists(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("error creating basepath: %d", err)
	}
	builder.path = opts.Path

	tilesPath := opts.Path + "/valhalla_tiles"
	err = createPathIfNotExists(opts.Path + "/valhalla_tiles")
	if err != nil {
		return nil, fmt.Errorf("error creating valhalla_tiles path: %s", err)
	}
	builder.tilesPath = tilesPath

	builder.extractPath = builder.path + "/valhalla_tiles.tar"
	builder.adminPath = builder.tilesPath + "/admin.sqlite"
	builder.configPath = builder.path + "/valhalla.json"

	builder.datasetPath = builder.path + "/" + opts.Dataset
	if _, err := os.Stat(builder.datasetPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dataset: %s, does not exist", builder.datasetPath)
	}

	builder.concurrency = opts.Concurrency
	builder.maxCacheSize = opts.MaxCacheSize

	return builder, nil
}

func (ve *TileBuilder) GetConfigPath() string {
	return ve.configPath
}

func (ve *TileBuilder) GetTilePath() string {
	return ve.tilesPath
}

func (ve *TileBuilder) GetExtractPath() string {
	return ve.extractPath
}

func (ve *TileBuilder) GetAdminPath() string {
	return ve.adminPath
}

func (ve *TileBuilder) GetDatasetPath() string {
	return ve.adminPath
}

func (ve *TileBuilder) GetConcurrency() int {
	return ve.concurrency
}

func (ve *TileBuilder) GetMaxCacheSize() int64 {
	return ve.maxCacheSize
}

func (ve *TileBuilder) BuildConfig(ctx context.Context) error {
	start := time.Now()
	args := []string{
		"--mjolnir-concurrency", fmt.Sprint(ve.concurrency),
		"--mjolnir-max-cache-size", fmt.Sprint(ve.maxCacheSize),
		"--mjolnir-tile-dir", ve.tilesPath,
		"--mjolnir-tile-extract", ve.extractPath,
		"--mjolnir-admin", ve.adminPath,
	}

	ve.logger.Info("creating valhalla config", "args", args)

	output, err := ve.executor.executeWithOutput(ctx, "valhalla_build_config", args)
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

	elapsed := time.Since(start)
	ve.logger.Info(
		"finished creating valhalla config",
		"args", args,
		"elapsed", elapsed.String(),
	)

	return nil
}

func (ve *TileBuilder) BuildTiles(ctx context.Context) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	start := time.Now()

	args := []string{"--config", ve.configPath, ve.datasetPath}
	ve.logger.Info("started creating tiles", "args", args)

	err := ve.executor.execute(ctx, "valhalla_build_tiles", args)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	ve.logger.Info(
		"finished creating tiles",
		"args", args,
		"elapsed", elapsed.String(),
	)

	return nil
}

func (ve *TileBuilder) BuildTilesExtract(ctx context.Context) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	start := time.Now()

	args := []string{"--config", ve.configPath, "-O"}
	ve.logger.Info("started tarballing tiles extract", "args", args)

	err := ve.executor.execute(ctx, "valhalla_build_extract", args)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	ve.logger.Info(
		"finished tarballing tile extract",
		"args", args,
		"elapsed", elapsed.String(),
	)

	return nil
}
