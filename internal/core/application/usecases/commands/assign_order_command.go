package commands

type AssignOrderCommand struct {
	isValid bool
}

func (c AssignOrderCommand) IsValid() bool {
	return c.isValid
}

func NewAssignOrderCommand() (AssignOrderCommand, error) {
	return AssignOrderCommand{isValid: true}, nil
}
