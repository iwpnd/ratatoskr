package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/iwpnd/valhalla-builder/services/tiles"

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

func (b *TestTileBuilder) BuildConfig(ctx context.Context) error {
	switch value := b.buildconfig.(type) {
	case error:
		return value
	}

	return nil
}

func TestConfigState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		params    Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildConfig returns an Error",
			params: Params{
				Logger: logger,
				Builder: &TestTileBuilder{
					buildconfig: fmt.Errorf("error during config creation"),
				},
			},
			expectErr: true,
		},
		{
			name: "configState returns buildState",
			params: Params{
				Logger:  logger,
				Builder: &TestTileBuilder{},
			},
			expectErr: false,
			wantState: buildState,
		},
	}

	for _, test := range tests {
		_, nextState, err := configState(ctx, test.params)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestConfigState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestConfigState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, gotState, wantState)
	}
}
