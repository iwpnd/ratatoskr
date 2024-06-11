package states

import (
	"context"
)

func compressState(ctx context.Context, args Args) (Args, State[Args], error) {
	if args.Compressor == nil {
		return args, nil, nil
	}

	args.Logger.Info("starting compression state", "name", args.Name)

	archive := args.Builder.GetPath() + "/valhalla_tiles"
	err := args.Compressor.Compress(
		ctx,
		archive,
		args.Builder.GetExtractPath(),
		args.Builder.GetAdminPath(),
	)

	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished compression state", "name", args.Name, "archive", archive)
	return args, nil, nil
}
