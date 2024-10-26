package states

import (
	"context"
	"time"
)

func DownloadState(ctx context.Context, params Params) (Params, State[Params], error) {
	if params.Downloader == nil {
		return params, ConfigState, nil
	}

	params.Logger.Info("starting download state", "name", params.Name)

	start := time.Now()

	err := params.Downloader.Get(ctx)
	if err != nil {
		return params, nil, &StateError{State: downloadState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished download state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, ConfigState, nil
}
