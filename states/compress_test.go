package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCompressor struct {
	compress any
}

func (c *TestCompressor) Compress(ctx context.Context, archive string, files ...string) error {
	switch value := c.compress.(type) { //nolint:gocritic
	case error:
		return value
	}

	return nil
}

func TestCompressState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := t.Context()

	tests := []struct {
		name      string
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "Compressor returns an Error",
			params: NewParams("test", "test", logger).
				WithTileBuilder(&TestTileBuilder{}).
				WithCompression(&TestCompressor{
					compress: fmt.Errorf("error during compression"),
				}),
			expectErr: true,
		},
		{
			name: "compressState returns nil",
			params: NewParams("test", "test", logger).
				WithTileBuilder(&TestTileBuilder{path: "testpath"}).
				WithCompression(&TestCompressor{}),
			expectErr: false,
			wantState: nil,
		},
		{
			name:      "No Compressor defined",
			params:    NewParams("test", "test", logger),
			expectErr: false,
			wantState: nil,
		},
	}

	for _, test := range tests {
		_, nextState, err := CompressState(ctx, test.params)
		switch {
		case err == nil && test.expectErr:
			t.Errorf("TestCompressState - %s: got err == nil, want err != nil", test.name)
			continue
		case err != nil && !test.expectErr:
			t.Errorf("TestCompressState - %s: got err != nil, want err == nil", test.name)
			continue
		case err != nil:
			continue
		}

		gotState := methodName(nextState)
		wantState := methodName(test.wantState)

		assert.Equal(t, wantState, gotState)
	}
}
