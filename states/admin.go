package states

import (
	"context"
	"time"
)

func adminState(ctx context.Context, params Params) (Params, State[Params], error) {
	params.Logger.Info("starting admins state", "name", params.Name)
	start := time.Now()

	err := params.Builder.BuildAdmins(ctx)
	if err != nil {
		return params, nil, &StateError{State: AdminState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished admins state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, extractState, nil
}
