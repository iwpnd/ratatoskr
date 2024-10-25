package states

import (
	"context"
	"time"
)

func configState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting config state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildConfig(ctx)
	if err != nil {
		return params, nil, &StateError{State: ConfigState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished config state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, buildState, nil
}
