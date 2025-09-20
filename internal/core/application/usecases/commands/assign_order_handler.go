package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type AssignOrdersCommandHandler interface {
	Handle(ctx context.Context, command AssignOrderCommand) error
}

var _ AssignOrdersCommandHandler = &assignOrderCommandHandler{}

type assignOrderCommandHandler struct {
	uowFactory      ports.UnitOfWorkFactory
	orderDispatcher services.OrderDispatcher
}

func NewAssignOrdersCommandHandler(uowFactory ports.UnitOfWorkFactory, orderDispatcher services.OrderDispatcher) (AssignOrdersCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}

	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return assignOrderCommandHandler{
		uowFactory:      uowFactory,
		orderDispatcher: orderDispatcher,
	}, nil
}

func (h assignOrderCommandHandler) Handle(ctx context.Context, command AssignOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("assign order command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	orderAggregate, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	couriers, err := uow.CourierRepository().GetAllFree(ctx)
	if err != nil {
		return err
	}

	if len(couriers) == 0 {
		return nil
	}

	courier, err := h.orderDispatcher.Dispatch(orderAggregate, couriers)
	if err != nil {
		return err
	}

	uow.Begin(ctx)

	err = uow.CourierRepository().Update(ctx, courier)
	if err != nil {
		return err
	}

	err = uow.OrderRepository().Update(ctx, orderAggregate)
	if err != nil {
		return err
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
