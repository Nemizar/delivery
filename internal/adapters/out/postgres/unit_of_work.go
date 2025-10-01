package postgres

import (
	"context"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	committed         bool
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
	mediatr           ddd.Mediatr
}

func NewUnitOfWork(db *gorm.DB, mediatr ddd.Mediatr) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	uow := &UnitOfWork{
		db:      db,
		mediatr: mediatr,
	}

	courierRepo, err := courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.courierRepository = courierRepo

	orderRepo, err := orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.orderRepository = orderRepo

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
	u.committed = false
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	if err := u.persistDomainEvents(ctx, u.tx); err != nil {
		return err
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	u.committed = true
	u.clearTx()
	return nil
}

func (u *UnitOfWork) RollbackUnlessCommitted(ctx context.Context) {
	if u.tx != nil && !u.committed {
		if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
			log.Error(err)
		}
		u.clearTx()
	}
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
	u.committed = false
}

func (u *UnitOfWork) persistDomainEvents(ctx context.Context, tx *gorm.DB) error {
	for _, aggregate := range u.trackedAggregates {
		outboxEvents, err := outbox.EncodeDomainEvents(aggregate.GetDomainEvents())
		if err != nil {
			return err
		}

		if len(outboxEvents) > 0 {
			if err := tx.WithContext(ctx).Create(&outboxEvents).Error; err != nil {
				return err
			}
		}

		aggregate.ClearDomainEvents()
	}

	return nil
}
