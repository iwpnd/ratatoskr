package states

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (b *TestTileBuilder) BuildAdmins(ctx context.Context) error {
	switch value := b.buildadmins.(type) {
	case error:
		return value
	}

	return nil
}

func TestAdminState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name      string
		args      Args
		expectErr bool
		wantState State[Args]
	}{
		{
			name: "TileBuilder.BuildAdmins returns an Error",
			args: Args{
				Logger: logger,
				Builder: &TestTileBuilder{
					buildadmins: fmt.Errorf("error during building admins"),
				},
			},
			expectErr: true,
		},
		{
			name: "adminState returns extractState",
			args: Args{
				Logger:  logger,
				Builder: &TestTileBuilder{},
			},
			expectErr: false,
			wantState: extractState,
		},
	}

	for _, test := range tests {
		_, nextState, err := adminState(ctx, test.args)
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

		assert.Equal(t, gotState, wantState)
	}
}
