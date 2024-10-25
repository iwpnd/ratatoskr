package states

import (
	"context"
	"time"
)

func buildState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting build state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildTiles(ctx)
	if err != nil {
		return params, nil, &StateError{State: BuildState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished build state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, adminState, nil
}
