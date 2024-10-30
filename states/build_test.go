package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildTiles returns an Error",
			params: NewParams("test", "test", logger).
				WithTileBuilder(
					&TestTileBuilder{
						buildtiles: fmt.Errorf("error during building tiles"),
					},
				),
			expectErr: true,
		},
		{
			name: "buildState returns adminState",
			params: NewParams("test", "test", logger).
				WithTileBuilder(
					&TestTileBuilder{},
				),
			expectErr: false,
			wantState: AdminState,
		},
	}

	for _, test := range tests {
		_, nextState, err := BuildState(ctx, test.params)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestBuildState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestBuildState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, gotState, wantState)
	}
}
