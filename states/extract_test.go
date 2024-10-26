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
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildTilesExtract returns an Error",
			params: NewParams("test", logger).
				WithTileBuilder(
					&TestTileBuilder{
						buildtilesextract: fmt.Errorf("error during building extract"),
					},
				),
			expectErr: true,
		},
		{
			name: "extractState returns compressState",
			params: NewParams("test", logger).
				WithTileBuilder(
					&TestTileBuilder{},
				),
			expectErr: false,
			wantState: CompressState,
		},
	}

	for _, test := range tests {
		_, nextState, err := ExtractState(ctx, test.params)
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
