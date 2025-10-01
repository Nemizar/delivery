package order

import (
	"delivery/internal/pkg/ddd"
	"reflect"

	"github.com/google/uuid"
)

var _ ddd.DomainEvent = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	ID   uuid.UUID
	Name string

	Payload DTO

	isValid bool
}

func (e *CompletedDomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *CompletedDomainEvent) GetName() string {
	return e.Name
}

func NewCompletedDomainEvent(aggregate *Order) *CompletedDomainEvent {
	dto := DTO{
		ID:     aggregate.ID(),
		Status: aggregate.Status(),
	}

	event := CompletedDomainEvent{
		ID:      uuid.New(),
		Payload: dto,
		isValid: true,
	}

	event.Name = reflect.TypeOf(event).Name()

	return &event
}

func NewEmptyCompletedDomainEvent() ddd.DomainEvent {
	event := CompletedDomainEvent{}
	event.Name = reflect.TypeOf(event).Name()
	return &event
}

func (e *CompletedDomainEvent) IsValid() bool {
	return e.isValid
}

type DTO struct {
	ID     uuid.UUID
	Status Status
}
