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
		params    *Params
		expectErr bool
		wantState State[Params]
	}{
		{
			name: "TileBuilder.BuildAdmins returns an Error",
			params: NewParams("test", logger).
				WithTileBuilder(
					&TestTileBuilder{buildadmins: fmt.Errorf("error during building admins")},
				),
			expectErr: true,
		},
		{
			name: "adminState returns extractState",
			params: NewParams("test", logger).
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

		assert.Equal(t, gotState, wantState)
	}
}
