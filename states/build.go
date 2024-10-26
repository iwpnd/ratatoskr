package states

import (
	"context"
	"time"
)

func BuildState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting build state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildTiles(ctx)
	if err != nil {
		return params, nil, &StateError{State: buildState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished build state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, AdminState, nil
}
