package states

import (
	"context"
	"time"
)

func ConfigState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting config state", "name", params.dataset)
	start := time.Now()

	err := params.builder.BuildConfig(ctx, params.dataset, params.outputPath)
	if err != nil {
		return params, nil, &StateError{State: configState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished config state",
		"name", params.dataset,
		"elapsed", elapsed.String(),
	)
	return params, BuildState, nil
}
