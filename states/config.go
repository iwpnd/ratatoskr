package states

import (
	"context"
)

func configState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting config state", "name", args.Name)

	err := args.Builder.BuildConfig(ctx)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished config state", "name", args.Name)
	return args, buildState, nil
}
