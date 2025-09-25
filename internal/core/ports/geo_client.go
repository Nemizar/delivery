package ports

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
)

type GeoClient interface {
	GetLocation(ctx context.Context, street string) (kernel.Location, error)
	Close() error
}
