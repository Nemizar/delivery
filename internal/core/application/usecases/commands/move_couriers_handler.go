package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCouriersCommandHandler interface {
	Handle(ctx context.Context, command MoveCouriersCommand) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCouriersCommandHandler(uowFactory ports.UnitOfWorkFactory) (MoveCouriersCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}

	return moveCouriersCommandHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h moveCouriersCommandHandler) Handle(ctx context.Context, command MoveCouriersCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("move couriers command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	assignedOrders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, order := range assignedOrders {
		uow.Begin(ctx)

		courier, err := uow.CourierRepository().Get(ctx, *order.CourierID())
		if err != nil {
			return err
		}

		err = courier.Move(order.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(order.Location()) {
			err = order.Complete()
			if err != nil {
				return err
			}

			err = courier.CompleteOrder(order)
			if err != nil {
				return err
			}
		}

		err = uow.OrderRepository().Update(ctx, order)
		if err != nil {
			return err
		}

		err = uow.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}

		err = uow.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
