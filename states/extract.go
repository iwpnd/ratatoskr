package states

import (
	"context"
	"time"
)

func extractState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting tiles extract state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildTilesExtract(ctx)
	if err != nil {
		return params, nil, &StateError{State: ExtractState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished tiles extract state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, compressState, nil
}
