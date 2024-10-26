package states

import (
	"context"
	"time"
)

func ConfigState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	if params.builder == nil {
		return params, nil, nil
	}

	params.logger.Info("starting config state", "name", params.name)
	start := time.Now()

	err := params.builder.BuildConfig(ctx)
	if err != nil {
		return params, nil, &StateError{State: configState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished config state",
		"name", params.name,
		"elapsed", elapsed.String(),
	)
	return params, BuildState, nil
}
