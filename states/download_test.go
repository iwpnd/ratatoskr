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

	"github.com/iwpnd/ratatoskr/services/download"

	"github.com/stretchr/testify/assert"
)

type TestDownloader struct {
	get any
	md5 any

	download.Downloader
}

func (td *TestDownloader) Get(ctx context.Context, dataset, outputPath string) error {
	switch value := td.get.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func (td *TestDownloader) MD5(ctx context.Context, dataset string) (string, error) {
	switch value := td.md5.(type) { //nolint:gocritic
	case error:
		return "", value
	}

	return td.md5.(string), nil //nolint:forcetypeassert
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
	ctx := t.Context()

	tests := []struct {
		name      string
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "Downloader returns an Error on Get",
			params: NewParams("test", "test", logger).
				WithDownload(
					&TestDownloader{
						get: fmt.Errorf("error during MD5"),
						md5: "md5",
					},
				),
			expectErr: true,
		},
		{
			name: "Downloader returns an Error on Get",
			params: NewParams("test", "test", logger).
				WithDownload(
					&TestDownloader{
						get: fmt.Errorf("error during Get"),
						md5: "md5",
					},
				),
			expectErr: true,
		},
		{
			name: "downloadState returns a configState",
			params: NewParams("test", "test", logger).
				WithDownload(
					&TestDownloader{
						get: nil,
						md5: "md5",
					},
				),
			expectErr: false,
			wantState: ConfigState,
		},
	}

	for _, test := range tests {
		_, nextState, err := DownloadState(ctx, test.params)
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

		assert.Equal(t, wantState, gotState)
	}
}
