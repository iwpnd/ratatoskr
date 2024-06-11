package states

import (
	"context"
	"time"
)

func buildState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting build state", "name", args.Name)
	start := time.Now()

	err := args.Builder.BuildTiles(ctx)
	if err != nil {
		return args, nil, err
	}

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished build state",
		"name", args.Name,
		"elapsed", elapsed.String(),
	)
	return args, adminState, nil
}
