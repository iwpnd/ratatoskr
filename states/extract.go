package states

import (
	"context"
	"time"
)

func ExtractState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting tiles extract state", "name", params.name)
	start := time.Now()

	err := params.builder.BuildTilesExtract(ctx)
	if err != nil {
		return params, nil, &StateError{State: extractState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished tiles extract state",
		"name", params.name,
		"elapsed", elapsed.String(),
	)
	return params, CompressState, nil
}
