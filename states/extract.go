package states

import (
	"context"
	"time"
)

func extractState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting tiles extract state", "name", args.Name)
	start := time.Now()

	err := args.Builder.BuildTilesExtract(ctx)
	if err != nil {
		return args, nil, err
	}

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished tiles extract state",
		"name", args.Name,
		"elapsed", elapsed.String(),
	)
	return args, compressState, nil
}
