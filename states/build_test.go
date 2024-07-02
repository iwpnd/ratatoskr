package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (b *TestTileBuilder) BuildTiles(ctx context.Context) error {
	switch value := b.buildtiles.(type) {
	case error:
		return value
	}

	return nil
}

func TestBuildState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		args      Args
		expectErr bool
		wantState State[Args]
	}{
		{
			name: "TileBuilder.BuildTiles returns an Error",
			args: Args{
				Logger: logger,
				Builder: &TestTileBuilder{
					buildtiles: fmt.Errorf("error during building tiles"),
				},
			},
			expectErr: true,
		},
		{
			name: "buildState returns adminState",
			args: Args{
				Logger:  logger,
				Builder: &TestTileBuilder{},
			},
			expectErr: false,
			wantState: adminState,
		},
	}

	for _, test := range tests {
		_, nextState, err := buildState(ctx, test.args)
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
