package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/iwpnd/valhalla-builder/services/download"

	"github.com/stretchr/testify/assert"
)

type TestDownloader struct {
	get any

	download.Downloader
}

func (td *TestDownloader) Get(ctx context.Context) error {
	switch value := td.get.(type) {
	case error:
		return value
	}

	return nil
}

func methodName(method any) string {
	if method == nil {
		return "<nil>"
	}
	return strings.TrimSuffix(
		runtime.FuncForPC(
			reflect.ValueOf(method).Pointer(),
		).Name(),
		"-fm")
}

func TestDownloadState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		args      Args
		expectErr bool
		wantState State[Args]
	}{
		{
			name: "Downloader returns an Error",
			args: Args{
				Logger: logger,
				Downloader: &TestDownloader{
					get: fmt.Errorf("error during download"),
				},
			},
			expectErr: true,
		},
		{
			name: "downloadState returns an configState",
			args: Args{
				Logger: logger,
				Downloader: &TestDownloader{
					get: nil,
				},
			},
			expectErr: false,
			wantState: configState,
		},
		{
			name: "No Downloader defined",
			args: Args{
				Logger: logger,
			},
			expectErr: false,
			wantState: configState,
		},
	}

	for _, test := range tests {
		_, nextState, err := downloadState(ctx, test.args)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestDownloadState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestDownloadState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, gotState, wantState)
	}
}
