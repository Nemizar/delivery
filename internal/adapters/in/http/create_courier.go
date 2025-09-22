package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) CreateCourier(ctx echo.Context) error {
	var c servers.NewCourier
	if err := ctx.Bind(&c); err != nil {
		return problems.NewBadRequest("invalid request body: " + err.Error())
	}

	command, err := commands.NewCreateCourierCommand(c.Name, c.Speed)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createCourierCommandHandler.Handle(ctx.Request().Context(), command)

	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}

		return problems.NewConflict(err.Error(), "/")
	}

	return ctx.JSON(http.StatusCreated, nil)
}
