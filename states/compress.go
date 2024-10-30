package states

import (
	"context"
	"fmt"
	"time"
)

func CompressState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	if params.compressor == nil {
		return params, nil, nil
	}

	params.logger.Info("starting compression state", "name", params.dataset)
	start := time.Now()

	path, ok := params.builder.Path()
	if !ok {
		return params, nil, &StateError{
			State: compressState,
			Err:   fmt.Errorf("builder path (%s) does not exist", path),
		}
	}
	tilesPath, ok := params.builder.TilesPath()
	if !ok {
		return params, nil, &StateError{
			State: compressState,
			Err:   fmt.Errorf("tiles path (%s) does not exist", tilesPath),
		}
	}
	extractPath, ok := params.builder.ExtractPath()
	if !ok {
		return params, nil, &StateError{
			State: compressState,
			Err:   fmt.Errorf("extract path (%s) does not exist", tilesPath),
		}
	}
	adminPath, ok := params.builder.AdminPath()
	if !ok {
		return params, nil, &StateError{
			State: compressState,
			Err:   fmt.Errorf("admin path (%s) does not exist", tilesPath),
		}
	}

	archive := path + "/valhalla_tiles"
	err := params.compressor.Compress(
		ctx,
		tilesPath,
		extractPath,
		adminPath,
	)

	if err != nil {
		return params, nil, &StateError{State: compressState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished compression state",
		"name",
		params.dataset,
		"archive", archive,
		"elapsed", elapsed.String(),
	)
	return params, nil, nil
}
