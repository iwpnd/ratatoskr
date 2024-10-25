package states

import (
	"context"
	"time"
)

func downloadState(ctx context.Context, params Params) (Params, State[Params], error) {
	if params.Downloader == nil {
		return params, configState, nil
	}

	params.Logger.Info("starting download state", "name", params.Name)

	start := time.Now()

	err := params.Downloader.Get(ctx)
	if err != nil {
		return params, nil, &StateError{State: DownloadState, Err: err}
	}

	elapsed := time.Since(start)
	params.Logger.Info(
		"successfully finished download state",
		"name", params.Name,
		"elapsed", elapsed.String(),
	)
	return params, configState, nil
}
