package states

import (
	"context"
	"time"
)

func downloadState(ctx context.Context, args Args) (Args, State[Args], error) {
	if args.Downloader == nil {
		return args, configState, nil
	}

	args.Logger.Info("starting download state", "name", args.Name)

	start := time.Now()

	err := args.Downloader.Get(ctx)
	if err != nil {
		args.Logger.Error("Failed download state", "name", args.Name, "error", err)
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
