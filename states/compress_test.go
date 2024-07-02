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
	switch value := c.compress.(type) {
	case error:
		return value
	}

	return nil

}

func (tb *TestTileBuilder) GetPath() string {
	return tb.path
}

func (tb *TestTileBuilder) GetExtractPath() string {
	return tb.path
}

func (tb *TestTileBuilder) GetAdminPath() string {
	return tb.path
}

func TestCompressState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		args      Args
		expectErr bool
		wantState State[Args]
	}{
		{
			name: "Compressor returns an Error",
			args: Args{
				Logger:  logger,
				Builder: &TestTileBuilder{path: "testpath"},
				Compressor: &TestCompressor{
					compress: fmt.Errorf("error during compression"),
				},
			},
			expectErr: true,
		},
		{
			name: "compressState returns nil",
			args: Args{
				Logger:     logger,
				Builder:    &TestTileBuilder{path: "testpath"},
				Compressor: &TestCompressor{},
			},
			expectErr: false,
			wantState: nil,
		},
		{
			name: "No Compressor defined",
			args: Args{
				Logger: logger,
			},
			expectErr: false,
			wantState: nil,
		},
	}

	for _, test := range tests {
		_, nextState, err := compressState(ctx, test.args)
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

		assert.Equal(t, gotState, wantState)
	}
}
