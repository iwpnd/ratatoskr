package states

import (
	"context"
)

type States int

const (
	adminState States = iota
	buildState
	compressState
	configState
	downloadState
	extractState
)

func (s States) String() string {
	switch s {
	case adminState:
		return "AdminState"
	case buildState:
		return "BuildState"
	case compressState:
		return "CompressState"
	case configState:
		return "ConfigState"
	case downloadState:
		return "DownloadState"
	case extractState:
		return "ExtractState"
	default:
		return "Unknown"
	}
}

type State[T any] func(ctx context.Context, params *T) (*T, State[T], error)
