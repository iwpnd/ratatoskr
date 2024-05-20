package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type TileBuilder struct {
	ConfigDir string

	logger   *slog.Logger
	executor *executor
}

type mjolnirConfig struct {
	TileDir string `json:"tile_dir"`
}

func NewValhallaExecutor(configDir string, logger *slog.Logger, debug bool) (*TileBuilder, error) {
	executor := &executor{logger: logger, debug: debug}
	if !executor.hasExecutable("valhalla_build_tiles") {
		return nil, fmt.Errorf("missing executable: valhalla_build_tiles")
	}

	if !executor.hasExecutable("valhalla_build_extract") {
		return nil, fmt.Errorf("missing executable: valhalla_build_extract")
	}

	if !executor.hasExecutable("valhalla_build_admins") {
		return nil, fmt.Errorf("missing executable: valhalla_build_admins")
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, err
	}

	configFile, err := os.Open(configDir)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	byteValue, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	configFileDir, _ := filepath.Split(configDir)

	var config map[string]mjolnirConfig
	err = json.Unmarshal([]byte(byteValue), &config)
	if err != nil {
		return nil, err
	}

	_, tilesDir := filepath.Split(config["mjolnir"].TileDir)
	if _, err := os.Stat(configFileDir + "/" + tilesDir); os.IsNotExist(err) {
		logger.Info("could not find tile directory: " + tilesDir + ", creating..")
		err = os.Mkdir(tilesDir, os.ModePerm)
		if err != nil {
			logger.Error("could not create tile directory: ", err)
			return nil, err
		}
		logger.Info("successfully created tile directory: " + tilesDir)
	}
	return &TileBuilder{executor: executor, logger: logger, ConfigDir: configDir}, nil
}

func (ve *TileBuilder) BuildTiles(ctx context.Context, args []string) error {
	start := time.Now()

	a := append([]string{"--config=" + ve.ConfigDir}, args...)
	ve.logger.Info("starting tile builder.", "args", a)

	err := ve.executor.execute(ctx, "valhalla_build_tiles", a)
	if err != nil {
		ve.logger.Error("cannot execute valhalla_build_tiles: ", "err", err)
		return err
	}

	elapsed := time.Since(start)
	ve.logger.Info("starting tile builder.", "args", a, "elapsed", elapsed.String())

	return nil
}
