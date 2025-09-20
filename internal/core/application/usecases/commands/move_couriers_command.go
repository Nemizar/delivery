package commands

type MoveCouriersCommand struct {
	isValid bool
}

func (c MoveCouriersCommand) IsValid() bool {
	return c.isValid
}

func NewMoveCouriersCommand() (MoveCouriersCommand, error) {
	return MoveCouriersCommand{isValid: true}, nil
}
