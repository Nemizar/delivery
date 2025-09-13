package services_test

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"delivery/internal/pkg/tests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderDispatcher_Dispatch_Success(t *testing.T) {
	couriers := []*courier.Courier{
		tests.CreateCourier("Bob", 1, tests.CreateLocation(1, 1)),
		tests.CreateCourier("Alice", 2, tests.CreateLocation(3, 3)),
	}

	o, err := order.NewOrder(uuid.New(), tests.CreateLocation(1, 1), 5)
	require.NoError(t, err)

	dispatcher := services.NewOrderDispatcher()

	got, err := dispatcher.Dispatch(o, couriers)
	require.NoError(t, err)
	assert.Equal(t, couriers[0], got)

	assert.Equal(t, order.StatusAssigned, o.Status())
	require.NotNil(t, o.CourierID())
	assert.Equal(t, couriers[0].Id(), *o.CourierID())
}

func TestOrderDispatcher_Dispatch_Errors(t *testing.T) {
	testCases := []struct {
		name     string
		order    *order.Order
		couriers []*courier.Courier
	}{
		{
			name:     "nil order",
			order:    nil,
			couriers: []*courier.Courier{},
		},
		{
			name:     "nil couriers slice",
			order:    tests.CreateOrder(uuid.New(), tests.CreateLocation(1, 1), 5),
			couriers: nil,
		},
		{
			name:     "empty couriers slice",
			order:    tests.CreateOrder(uuid.New(), tests.CreateLocation(1, 1), 5),
			couriers: []*courier.Courier{},
		},
		{
			name:     "no suitable courier (capacity too low)",
			order:    tests.CreateOrder(uuid.New(), tests.CreateLocation(1, 1), 15),
			couriers: []*courier.Courier{tests.CreateCourier("Bob", 1, tests.CreateLocation(1, 1))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dispatcher := services.NewOrderDispatcher()
			got, err := dispatcher.Dispatch(tc.order, tc.couriers)

			assert.Nil(t, got)
			assert.Error(t, err)
		})
	}
}

func TestOrderDispatcher_Dispatch_OrderCompleted(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	o := testsCreateOrderWithLocationAndWeight(orderID, 1, 1, 5)

	c := tests.CreateCourier("Bob", 10, tests.CreateLocation(1, 1))
	err := o.Assign(c.Id())
	require.NoError(t, err)
	err = o.Complete()
	require.NoError(t, err)

	dispatcher := services.NewOrderDispatcher()
	got, err := dispatcher.Dispatch(o, []*courier.Courier{c})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, services.ErrOrderIsAlreadyAssigned)
}

func testsCreateOrderWithLocationAndWeight(id uuid.UUID, x, y, weight int) *order.Order {
	loc := tests.CreateLocation(x, y)
	o, err := order.NewOrder(id, loc, weight)
	if err != nil {
		panic(err)
	}
	return o
}

func Test_OrderDispatcher_Errors(t *testing.T) {
	od := services.NewOrderDispatcher()
	c := tests.CreateCourier("Bob", 1, tests.CreateLocation(1, 1))
	o := tests.CreateOrder(uuid.New(), tests.CreateLocation(1, 1), 15)
	err := o.Assign(tests.CreateCourier("Bob", 1, tests.CreateLocation(1, 1)).Id())
	err = o.Complete()
	assert.Nil(t, err)
	c, err = od.Dispatch(o, []*courier.Courier{
		tests.CreateCourier("Bob", 1, tests.CreateLocation(1, 1)),
		tests.CreateCourier("Alice", 2, tests.CreateLocation(3, 3)),
	})
	assert.Nil(t, c)
	assert.ErrorIs(t, services.ErrOrderIsAlreadyAssigned, err)
}
