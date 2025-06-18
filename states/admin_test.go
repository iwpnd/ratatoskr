package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/iwpnd/ratatoskr/services/tiles"
	"github.com/stretchr/testify/assert"
)

type TestTileBuilder struct {
	buildconfig       any
	buildtiles        any
	buildadmins       any
	buildtilesextract any

	path string

	tiles.Builder
}

func (tb *TestTileBuilder) BuildConfig(
	ctx context.Context,
	dataset string,
	outputPath string,
) error {
	switch value := tb.buildconfig.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func (tb *TestTileBuilder) BuildAdmins(
	ctx context.Context,
	dataset string,
	outputPath string,
) error {
	switch value := tb.buildadmins.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func (tb *TestTileBuilder) BuildTilesExtract(
	ctx context.Context,
	dataset string,
	outputPath string,
) error {
	switch value := tb.buildtilesextract.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func (tb *TestTileBuilder) BuildTiles(
	ctx context.Context,
	dataset string,
	outputPath string,
) error {
	switch value := tb.buildtiles.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func (tb *TestTileBuilder) Path() (string, bool) {
	return tb.path, true
}

func (tb *TestTileBuilder) ExtractPath() (string, bool) {
	return tb.path, true
}

func (tb *TestTileBuilder) AdminPath() (string, bool) {
	return tb.path, true
}

func (tb *TestTileBuilder) TilesPath() (string, bool) {
	return tb.path, true
}

func TestAdminState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := t.Context()

	tests := []struct {
		name      string
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildAdmins returns an Error",
			params: NewParams("test", "test", logger).
				WithTileBuilder(
					&TestTileBuilder{buildadmins: fmt.Errorf("error during building admins")},
				),
			expectErr: true,
		},
		{
			name: "adminState returns extractState",
			params: NewParams("test", "test", logger).
				WithTileBuilder(
					&TestTileBuilder{},
				),
			expectErr: false,
			wantState: ExtractState,
		},
	}

	for _, test := range tests {
		_, nextState, err := AdminState(ctx, test.params)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestAdminState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestAdminState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, wantState, gotState)
	}
}
