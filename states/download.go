package states

import (
	"context"
)

func downloadState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting download state", "name", args.Name)

	err := args.Downloader.Get(ctx)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished download state", "name", args.Name)
	return args, configState, nil
}
