package states

import (
	"context"
)

type States int

const (
	AdminState States = iota
	BuildState
	CompressState
	ConfigState
	DownloadState
	ExtractState
)

func (s States) String() string {
	switch s {
	case AdminState:
		return "AdminState"
	case BuildState:
		return "BuildState"
	case CompressState:
		return "CompressState"
	case ConfigState:
		return "ConfigState"
	case DownloadState:
		return "DownloadState"
	case ExtractState:
		return "ExtractState"
	default:
		return "Unknown"
	}
}

type State[T any] func(ctx context.Context, params T) (T, State[T], error)

func Execute(ctx context.Context, params Params) error {
	if err := params.validate(ctx); err != nil {
		return err
	}

	start := downloadState
	_, err := run(ctx, params, start)
	if err != nil {
		return err
	}
	return nil
}

func run[T any](ctx context.Context, params T, start State[T]) (T, error) {
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
