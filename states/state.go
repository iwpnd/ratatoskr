package states

import (
	"context"
)

// States enum.
type States int

const (
	adminState States = iota
	buildState
	compressState
	configState
	downloadState
	extractState
)

// String to resolve the actual state name
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

// State s in the pipeline state machine
type State[T any] func(ctx context.Context, params *T) (*T, State[T], error)
