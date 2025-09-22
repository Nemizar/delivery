package commands

import "delivery/internal/pkg/errs"

type CreateCourierCommand struct {
	name    string
	speed   int
	isValid bool
}

func (c CreateCourierCommand) IsValid() bool {
	return c.isValid
}

func (c CreateCourierCommand) Speed() int {
	return c.speed
}

func (c CreateCourierCommand) Name() string {
	return c.name
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {
	if name == "" {
		return CreateCourierCommand{}, errs.NewValueIsInvalidError("name")
	}

	if speed <= 0 {
		return CreateCourierCommand{}, errs.NewValueIsInvalidError("speed")
	}

	return CreateCourierCommand{
		name:    name,
		speed:   speed,
		isValid: true,
	}, nil
}
