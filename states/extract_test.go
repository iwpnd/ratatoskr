package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (b *TestTileBuilder) BuildTilesExtract(ctx context.Context) error {
	switch value := b.buildtilesextract.(type) {
	case error:
		return value
	}

	return nil
}

func TestExtractState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		args      Args
		expectErr bool
		wantState State[Args]
	}{
		{
			name: "TileBuilder.BuildTilesExtract returns an Error",
			args: Args{
				Logger: logger,
				Builder: &TestTileBuilder{
					buildtilesextract: fmt.Errorf("error during building extract"),
				},
			},
			expectErr: true,
		},
		{
			name: "extractState returns compressState",
			args: Args{
				Logger:  logger,
				Builder: &TestTileBuilder{},
			},
			expectErr: false,
			wantState: compressState,
		},
	}

	for _, test := range tests {
		_, nextState, err := extractState(ctx, test.args)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestExtractState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestExtractState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, gotState, wantState)
	}
}
