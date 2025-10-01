package postgres

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type unitOfWorkFactory struct {
	db      *gorm.DB
	mediatr ddd.Mediatr
}

func NewUnitOfWorkFactory(db *gorm.DB, mediatr ddd.Mediatr) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &unitOfWorkFactory{db: db, mediatr: mediatr}, nil
}

func (f *unitOfWorkFactory) New(ctx context.Context) (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db.WithContext(ctx), f.mediatr)
}
