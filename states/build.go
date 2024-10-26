package states

import (
	"context"
	"time"
)

func BuildState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting build state", "name", params.name)
	start := time.Now()

	err := params.builder.BuildTiles(ctx)
	if err != nil {
		return params, nil, &StateError{State: buildState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished build state",
		"name", params.name,
		"elapsed", elapsed.String(),
	)
	return params, AdminState, nil
}
