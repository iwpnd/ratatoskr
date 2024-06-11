package states

import (
	"context"
)

func extractState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting tiles extract state", "name", args.Name)

	err := args.Builder.BuildTilesExtract(ctx)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished tiles extract state", "name", args.Name)
	return args, compressState, nil
}
