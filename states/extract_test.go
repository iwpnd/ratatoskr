package states

import (
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := t.Context()

	tests := []struct {
		name      string
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildTilesExtract returns an Error",
			params: NewParams("test", "test", logger).
				WithTileBuilder(
					&TestTileBuilder{
						buildtilesextract: fmt.Errorf("error during building extract"),
					},
				),
			expectErr: true,
		},
		{
			name: "extractState returns compressState",
			params: NewParams("test", "test", logger).
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

		assert.Equal(t, wantState, gotState)
	}
}
