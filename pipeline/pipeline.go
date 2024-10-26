package pipeline

import (
	"context"

	"github.com/iwpnd/valhalla-tiles-builder/states"
)

func Execute(ctx context.Context, params states.Params) error {
	if err := params.Validate(ctx); err != nil {
		return err
	}

	start := states.DownloadState
	_, err := run(ctx, params, start)
	if err != nil {
		return err
	}
	return nil
}

func run[T any](ctx context.Context, params T, start states.State[T]) (T, error) {
	var err error
	current := start
	for {
		if ctx.Err() != nil {
			return params, ctx.Err()
		}
		params, current, err = current(ctx, params)
		if err != nil {
			return params, err
		}
		if current == nil {
			return params, nil
		}
	}
}
