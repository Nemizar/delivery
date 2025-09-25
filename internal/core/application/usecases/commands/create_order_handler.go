package commands

import (
	"context"
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
	geoClient  ports.GeoClient
}

func NewCreateOrderCommandHandler(uowFactory ports.UnitOfWorkFactory, geoClient ports.GeoClient) (CreateOrderCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}

	if geoClient == nil {
		return nil, errs.NewValueIsRequiredError("geoClient")
	}

	return createOrderCommandHandler{
		uowFactory: uowFactory,
		geoClient:  geoClient,
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

	l, err := h.geoClient.GetLocation(ctx, command.Street())
	if err != nil {
		return err
	}

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
