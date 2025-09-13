package orderrepo

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.OrderRepository = &Repository{}

type Repository struct {
	tracker Tracker
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	return r.withTx(ctx, func(tx *gorm.DB) error {
		dto := DomainToDTO(aggregate)

		err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	return r.withTx(ctx, func(tx *gorm.DB) error {
		dto := DomainToDTO(aggregate)

		err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, id)
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Order", id)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	dto := OrderDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusCreated).
		First(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Created order", nil)
		}
		return nil, result.Error
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusAssigned).
		Find(&dtos)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return []*order.Order{}, nil
	}

	aggregates := make([]*order.Order, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}

func (r *Repository) withTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	isInTx := r.tracker.InTx()
	if !isInTx {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	if err := fn(tx); err != nil {
		return err
	}

	if !isInTx {
		return r.tracker.Commit(ctx)
	}
	return nil
}
