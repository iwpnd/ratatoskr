package states

import (
	"context"
	"time"
)

func configState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting config state", "name", args.Name)
	start := time.Now()

	err := args.Builder.BuildConfig(ctx)
	if err != nil {
		return args, nil, err
	}

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished config state",
		"name", args.Name,
		"elapsed", elapsed.String(),
	)
	return args, buildState, nil
}
