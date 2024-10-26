package states

import (
	"context"
	"time"
)

func AdminState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting admins state", "name", params.name)
	start := time.Now()

	err := params.builder.BuildAdmins(ctx)
	if err != nil {
		return params, nil, &StateError{State: adminState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished admins state",
		"name", params.name,
		"elapsed", elapsed.String(),
	)
	return params, ExtractState, nil
}
