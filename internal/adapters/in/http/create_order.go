package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s Server) CreateOrder(ctx echo.Context) error {
	createOrderCommand, err := commands.NewCreateOrderCommand(uuid.New(), "Несуществующая", 5)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrderCommandHandler.Handle(ctx.Request().Context(), createOrderCommand)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}

		return problems.NewConflict(err.Error(), "/")
	}

	return ctx.JSON(http.StatusCreated, nil)
}
