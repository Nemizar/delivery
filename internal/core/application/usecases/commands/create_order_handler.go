package commands

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type CreateOrderCommandHandler interface {
	Handle(ctx context.Context, command CreateOrderCommand) error
}

var _ CreateOrderCommandHandler = &createOrderCommandHandler{}

type createOrderCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateOrderCommandHandler(uowFactory ports.UnitOfWorkFactory) (CreateOrderCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}

	return createOrderCommandHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h createOrderCommandHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("create order command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	orderAggregate, err := uow.OrderRepository().Get(ctx, command.OrderID())
	if err != nil {
		if !errors.Is(err, errs.ErrObjectNotFound) {
			return err
		}
	}

	if orderAggregate != nil {
		return nil
	}

	l := kernel.RandomLocation()

	orderAggregate, err = order.NewOrder(command.OrderID(), l, command.Volume())
	if err != nil {
		return err
	}

	err = uow.OrderRepository().Add(ctx, orderAggregate)
	if err != nil {
		return err
	}

	return nil
}
