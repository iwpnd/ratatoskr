package states

import (
	"context"
)

func buildState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting build state", "name", args.Name)

	err := args.Builder.BuildTiles(ctx)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished build state", "name", args.Name)
	return args, adminState, nil
}
