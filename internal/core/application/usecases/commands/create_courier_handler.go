package commands

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(ctx context.Context, command CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateCourierCommandHandler(uowFactory ports.UnitOfWorkFactory) (CreateCourierCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}

	return createCourierCommandHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h createCourierCommandHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("create courier command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	l := kernel.RandomLocation()

	courierAggregate, err := courier.NewCourier(command.Name(), command.Speed(), l)
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Add(ctx, courierAggregate)
	if err != nil {
		return err
	}

	return nil
}
