package states

import (
	"context"
	"time"
)

func DownloadState(ctx context.Context, params *Params) (*Params, State[Params], error) {
	params.logger.Info("starting download state", "name", params.dataset)

	start := time.Now()

	md5, err := params.downloader.MD5(ctx, params.dataset)
	if err != nil {
		return params, nil, &StateError{State: downloadState, Err: err}
	}

	err = params.setOutputPath(params.dataset, md5)
	if err != nil {
		return params, nil, &StateError{State: downloadState, Err: err}
	}

	err = params.downloader.Get(ctx, params.dataset, params.outputPath)
	if err != nil {
		return params, nil, &StateError{State: downloadState, Err: err}
	}

	elapsed := time.Since(start)
	params.logger.Info(
		"successfully finished download state",
		"name", params.dataset,
		"elapsed", elapsed.String(),
	)
	return params, ConfigState, nil
}
