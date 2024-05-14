package states

import (
	"context"
)

func endState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting end state", "name", args.Name)
	err := args.BlobStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}

	err = args.FileStore.Get(ctx, args.Name)
	if err != nil {
		return args, nil, err
	}
	args.Logger.Info("finished end state", "name", args.Name)
	return args, nil, nil
}
