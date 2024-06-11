package states

import (
	"context"
)

func adminState(ctx context.Context, args Args) (Args, State[Args], error) {
	args.Logger.Info("starting admins state", "name", args.Name)

	err := args.Builder.BuildAdmins(ctx)
	if err != nil {
		return args, nil, err
	}

	args.Logger.Info("successfully finished admins state", "name", args.Name)
	return args, extractState, nil
}
