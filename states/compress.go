package states

import (
	"context"
	"time"
)

func CompressState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	if params.compressor == nil {
		return params, nil, nil
	}

	params.logger.Info("starting compression state", "name", params.name)
	start := time.Now()

	archive := params.builder.Path() + "/valhalla_tiles"
	err := params.compressor.Compress(
		ctx,
		params.builder.TilesPath(),
		params.builder.ExtractPath(),
		params.builder.AdminPath(),
	)

	if err != nil {
		return params, nil, &StateError{State: compressState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished compression state",
		"name",
		params.name,
		"archive", archive,
		"elapsed", elapsed.String(),
	)
	return params, nil, nil
}
