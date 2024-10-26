package states

import (
	"context"
	"time"
)

func AdminState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting admins state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildAdmins(ctx)
	if err != nil {
		return params, nil, &StateError{State: adminState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished admins state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, ExtractState, nil
}
