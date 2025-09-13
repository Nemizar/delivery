package order

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrderStatus     = errors.New("invalid order status")
	ErrCourierAlreadyAssigned = errors.New("courier already assigned")
	ErrOrderAlreadyCompleted  = errors.New("order already completed")
	ErrOrderNotAssigned       = errors.New("order not assigned")
)

type Order struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	courierID     *uuid.UUID
	location      kernel.Location
	volume        int
	status        Status
}

func NewOrder(id uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsInvalidError("id")
	}

	if volume <= 0 {
		return nil, errs.NewValueIsInvalidError("volume")
	}

	if !location.IsValid() {
		return nil, errs.NewValueIsInvalidError("location")
	}

	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		location:      location,
		volume:        volume,
		status:        StatusCreated,
	}, nil
}

func RestoreOrder(id uuid.UUID, courierID *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		courierID:     courierID,
		location:      location,
		volume:        volume,
		status:        status,
	}
}

func (o *Order) ID() uuid.UUID {
	return o.baseAggregate.ID()
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) ClearDomainEvents() {
	o.baseAggregate.ClearDomainEvents()
}

func (o *Order) GetDomainEvents() []ddd.DomainEvent {
	return o.baseAggregate.GetDomainEvents()
}

func (o *Order) RaiseDomainEvent(event ddd.DomainEvent) {
	o.baseAggregate.RaiseDomainEvent(event)
}

func (o *Order) Equals(other *Order) bool {
	return o.ID() == other.ID()
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o.courierID != nil {
		return ErrCourierAlreadyAssigned
	}

	if courierID == uuid.Nil {
		return errs.NewValueIsInvalidError("courierID")
	}

	if o.status != StatusCreated {
		return ErrInvalidOrderStatus
	}

	o.courierID = &courierID
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o.status == StatusCompleted {
		return ErrOrderAlreadyCompleted
	}

	if o.status != StatusAssigned {
		return ErrOrderNotAssigned
	}

	if o.courierID == nil {
		return errs.NewValueIsInvalidError("courierID")
	}

	o.status = StatusCompleted

	return nil
}
