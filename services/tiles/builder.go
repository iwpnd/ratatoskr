package tiles

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Builder interface {
	BuildConfig(ctx context.Context, dataset string, outputPath string) error
	BuildTiles(ctx context.Context, dataset string, outputPath string) error
	BuildTilesExtract(ctx context.Context, dataset string, outputPath string) error
	BuildAdmins(ctx context.Context, dataset string, outputPath string) error
	Path() (string, bool)
	AdminPath() (string, bool)
	ExtractPath() (string, bool)
	TilesPath() (string, bool)
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
	logger *slog.Logger) (*TileBuilder, error) {

	opts := &TileBuilderOptions{}

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

	return builder, nil
}

func (ve *TileBuilder) prepareWorkspace(_ context.Context, dataset string, outputPath string) error {
	if dataset == "" {
		return fmt.Errorf("missing dataset")
	}

	if outputPath == "" {
		return fmt.Errorf("missing outputPath")
	}

	err := createPathIfNotExists(outputPath)
	if err != nil {
		return fmt.Errorf("error creating basepath: %d", err)
	}
	ve.path = outputPath
	ve.tilesPath = outputPath + "/valhalla_tiles"
	err = createPathIfNotExists(ve.tilesPath)
	if err != nil {
		return fmt.Errorf("error creating valhalla_tiles path: %s", err)
	}
	ve.extractPath = ve.path + "/valhalla_tiles.tar"
	ve.adminPath = ve.path + "/admin.sqlite"
	ve.configPath = ve.path + "/config.json"
	ve.datasetPath = ve.path + "/" + toDatasetFileName(dataset)

	return nil
}

func (ve *TileBuilder) BuildConfig(ctx context.Context, dataset string, outputPath string) error {
	err := ve.prepareWorkspace(ctx, dataset, outputPath)
	if err != nil {
		return err
	}

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

func (ve *TileBuilder) BuildTiles(ctx context.Context, dataset string, outputPath string) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	err := ve.prepareWorkspace(ctx, dataset, outputPath)
	if err != nil {
		return err
	}

	params := []string{"--config", ve.configPath, ve.datasetPath}
	ve.logger.Info("started creating tiles", "params", params)

	err = ve.executor.execute(ctx, "valhalla_build_tiles", params)
	if err != nil {
		return err
	}

	ve.logger.Info(
		"finished creating tiles",
		"params", params,
	)

	return nil
}

func (ve *TileBuilder) BuildTilesExtract(ctx context.Context, dataset string, outputPath string) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	err := ve.prepareWorkspace(ctx, dataset, outputPath)
	if err != nil {
		return err
	}

	start := time.Now()

	params := []string{"--config", ve.configPath, "-O"}
	ve.logger.Info("started tarballing tiles extract", "params", params)

	err = ve.executor.execute(ctx, "valhalla_build_extract", params)
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

func (ve *TileBuilder) BuildAdmins(ctx context.Context, dataset string, outputPath string) error {
	if !ve.configCreated {
		return fmt.Errorf("error, create config first")
	}

	err := ve.prepareWorkspace(ctx, dataset, outputPath)
	if err != nil {
		return err
	}

	params := []string{"--config", ve.configPath, ve.datasetPath}
	ve.logger.Info("started building admins", "params", params)

	err = ve.executor.execute(ctx, "valhalla_build_admins", params)
	if err != nil {
		return err
	}

	ve.logger.Info(
		"finished building admins",
		"params", params,
	)

	return nil
}

func (ve *TileBuilder) Path() (string, bool) {
	if ve.path == "" {
		return "", false
	}
	return ve.path, true
}

func (ve *TileBuilder) ExtractPath() (string, bool) {
	if ve.extractPath == "" {
		return "", false
	}
	return ve.extractPath, true
}

func (ve *TileBuilder) AdminPath() (string, bool) {
	if ve.adminPath == "" {
		return "", false
	}
	return ve.adminPath, true
}

func (ve *TileBuilder) TilesPath() (string, bool) {
	if ve.tilesPath == "" {
		return "", false
	}
	return ve.tilesPath, true
}

func toDatasetFileName(dataset string) string {
	if dataset == "" {
		return ""
	}

	parts := strings.Split(dataset, "/")
	return parts[len(parts)-1] + "-latest.osm.pbf"
}
