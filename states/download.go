package states

import (
	"context"
	"time"
)

func DownloadState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	if params.downloader == nil {
		return params, ConfigState, nil
	}

	params.logger.Info("starting download state", "name", params.name)

	start := time.Now()

	err := params.downloader.Get(ctx)
	if err != nil {
		return params, nil, &StateError{State: downloadState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished download state",
		"name", params.name,
		"elapsed", elapsed.String(),
	)
	return params, ConfigState, nil
}
