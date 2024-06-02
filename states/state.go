package states

import (
	"context"
)

func Execute(ctx context.Context, args Args) error {
	if err := args.validate(ctx); err != nil {
		return err
	}

	start := downloadState
	_, err := run(ctx, args, start)
	if err != nil {
		return err
	}
	return nil
}

type State[T any] func(ctx context.Context, args T) (T, State[T], error)

func run[T any](ctx context.Context, args T, start State[T]) (T, error) {
	var err error
	current := start
	for {
		if ctx.Err() != nil {
			return args, ctx.Err()
		}
		args, current, err = current(ctx, args)
		if err != nil {
			return args, err
		}
		if current == nil {
			return args, nil
		}
	}
}
