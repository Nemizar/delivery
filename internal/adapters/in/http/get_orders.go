package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) GetOrders(ctx echo.Context) error {
	query, err := queries.NewGetNotCompletedOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getNotCompletedOrdersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
	}

	var httpResponse = make([]servers.Order, 0, len(queryResponse.Orders))
	for _, courier := range queryResponse.Orders {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var sCourier = servers.Order{
			Id:       courier.ID,
			Location: location,
		}

		httpResponse = append(httpResponse, sCourier)
	}

	return ctx.JSON(http.StatusOK, httpResponse)
}
