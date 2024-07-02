package states

import (
	"context"
	"time"
)

func downloadState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting download state", "name", args.Name)
	if args.Downloader == nil {
		return args, configState, nil
	}

	start := time.Now()

	err := args.Downloader.Get(ctx)
	if err != nil {
		return args, nil, err
	}

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished download state",
		"name", args.Name,
		"elapsed", elapsed.String(),
	)
	return args, configState, nil
}
