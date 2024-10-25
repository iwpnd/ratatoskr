package stores

import "context"

type Storer interface {
	Get(ctx context.Context, id string) error
}
