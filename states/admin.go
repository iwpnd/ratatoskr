package states

import (
	"context"
	"time"
)

func adminState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting admins state", "name", args.Name)
	start := time.Now()

	err := args.Builder.BuildAdmins(ctx)
	if err != nil {
		return args, nil, err
	}

	elapsed := time.Since(start)
	args.Logger.Info(
		"successfully finished admins state",
		"name", args.Name,
		"elapsed", elapsed.String(),
	)
	return args, extractState, nil
}
