package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/pkg/errs"
)

type Server struct {
	createCourierCommandHandler commands.CreateCourierCommandHandler
	createOrderCommandHandler   commands.CreateOrderCommandHandler

	getAllCouriersQueryHandler        queries.GetAllCouriersQueryHandler
	getNotCompletedOrdersQueryHandler queries.GetNotCompletedOrdersQueryHandler
}

func NewServer(
	createCourierCommandHandler commands.CreateCourierCommandHandler,
	createOrderCommandHandler commands.CreateOrderCommandHandler,
	getAllCouriersQueryHandler queries.GetAllCouriersQueryHandler,
	getNotCompletedOrdersQueryHandler queries.GetNotCompletedOrdersQueryHandler,
) (*Server, error) {
	if createCourierCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createCourierCommandHandler")
	}

	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	if getAllCouriersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersQueryHandler")
	}

	if getNotCompletedOrdersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getNotCompletedOrdersQueryHandler")
	}

	return &Server{
		createCourierCommandHandler:       createCourierCommandHandler,
		createOrderCommandHandler:         createOrderCommandHandler,
		getAllCouriersQueryHandler:        getAllCouriersQueryHandler,
		getNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
	}, nil
}
