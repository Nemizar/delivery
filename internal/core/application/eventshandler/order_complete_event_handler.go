package eventshandler

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

type orderCompletedEventHandler struct {
	orderProducer ports.OrderProducer
}

func NewOrderCompletedEventHandler(orderProducer ports.OrderProducer) (ddd.EventHandler, error) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("orderProducer")
	}

	return &orderCompletedEventHandler{
		orderProducer: orderProducer,
	}, nil
}

func (h *orderCompletedEventHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	err := h.orderProducer.Publish(ctx, domainEvent)
	if err != nil {
		return err
	}

	return nil
}
