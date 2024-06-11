package states

import (
	"context"
	"time"
)

func compressState(ctx context.Context, args Args) (Args, State[Args], error) {
	if args.Compressor == nil {
		return args, nil, nil
	}

	args.Logger.Info("starting compression state", "name", args.Name)
	start := time.Now()

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

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished compression state",
		"name",
		args.Name,
		"archive", archive,
		"elapsed", elapsed.String(),
	)
	return args, nil, nil
}
