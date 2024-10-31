package states

import (
	"context"
	"time"
)

// BuildState to build valhalla routing tiles from OSM dataset
func BuildState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting build state", "name", params.dataset)
	start := time.Now()

	err := params.builder.BuildTiles(ctx, params.dataset, params.outputPath)
	if err != nil {
		return params, nil, &StateError{State: buildState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished build state",
		"name", params.dataset,
		"elapsed", elapsed.String(),
	)
	return params, AdminState, nil
}
