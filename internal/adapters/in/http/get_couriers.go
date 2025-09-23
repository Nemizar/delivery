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

func (s Server) GetCouriers(ctx echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getAllCouriersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
	}

	var httpResponse = make([]servers.Courier, 0, len(queryResponse.Couriers))
	for _, courier := range queryResponse.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var sCourier = servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}

		httpResponse = append(httpResponse, sCourier)
	}

	return ctx.JSON(http.StatusOK, httpResponse)
}
