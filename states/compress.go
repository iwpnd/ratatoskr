package states

import (
	"context"
	"time"
)

func compressState(ctx context.Context, params Params) (Params, State[Params], error) {
	if params.Compressor == nil {
		return params, nil, nil
	}

	params.Logger.Info("starting compression state", "name", params.Name)
	start := time.Now()

	archive := params.Builder.GetPath() + "/valhalla_tiles"
	err := params.Compressor.Compress(
		ctx,
		archive,
		params.Builder.GetExtractPath(),
		params.Builder.GetAdminPath(),
	)

	if err != nil {
		return params, nil, &StateError{State: CompressState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished compression state",
		"name",
		params.Name,
		"archive", archive,
		"elapsed", elapsed.String(),
	)
	return params, nil, nil
}
