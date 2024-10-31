package states

import (
	"context"
	"time"
)

// AdminState to build admins from OSM dataset
func AdminState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting admins state", "name", params.dataset)
	start := time.Now()

	err := params.builder.BuildAdmins(ctx, params.dataset, params.outputPath)
	if err != nil {
		return params, nil, &StateError{State: adminState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished admins state",
		"name", params.dataset,
		"elapsed", elapsed.String(),
	)
	return params, ExtractState, nil
}
