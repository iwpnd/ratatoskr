package states

import (
	"context"
	"time"
)

// ExtractState to build a valhalla routing tiles .tar file
func ExtractState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting tiles extract state", "name", params.dataset)
	start := time.Now()

	err := params.builder.BuildTilesExtract(ctx, params.dataset, params.outputPath)
	if err != nil {
		return params, nil, &StateError{State: extractState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished tiles extract state",
		"name", params.dataset,
		"elapsed", elapsed.String(),
	)
	return params, CompressState, nil
}
