package states

import (
	"context"
)

func midState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting mid state", "name", args.Name)
	err := args.BlobStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}

	err = args.FileStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}
	args.Logger.Info("finished mid state", "name", args.Name)
	return args, endState, nil
}
