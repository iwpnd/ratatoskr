package states

import (
	"context"
)

func startState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting start state", "name", args.Name)

	err := args.BlobStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}

	err = args.FileStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("finished starting state", "name", args.Name)
	return args, midState, nil
}
