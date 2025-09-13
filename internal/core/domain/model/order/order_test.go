package order_test

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	id := uuid.New()
	location, err := kernel.NewLocation(5, 5)
	require.Nil(t, err)

	tests := []struct {
		name     string
		orderID  uuid.UUID
		volume   int
		location kernel.Location
		wantErr  error
	}{
		{
			name:     "Valid order",
			orderID:  id,
			volume:   10,
			location: location,
		},
		{
			name:    "Zero volume",
			orderID: id,
			volume:  0,
			wantErr: errs.NewValueIsInvalidError("volume"),
		},
		{
			name:     "Invalid location",
			orderID:  id,
			volume:   1,
			location: kernel.Location{},
			wantErr:  errs.NewValueIsInvalidError("location"),
		},
		{
			name:     "Nil order id",
			orderID:  uuid.Nil,
			volume:   1,
			location: location,
			wantErr:  errs.NewValueIsInvalidError("id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.NewOrder(tt.orderID, tt.location, tt.volume)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)

				return
			}

			assert.Equal(t, id, got.ID())
			assert.Equal(t, location, got.Location())
			assert.Equal(t, tt.volume, got.Volume())
			assert.Equal(t, order.StatusCreated, got.Status())
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	eqOrder := createValidOrder(kernel.RandomLocation(), 10)

	tests := []struct {
		name   string
		first  *order.Order
		second *order.Order
		want   bool
	}{
		{
			name:   "Equals",
			first:  eqOrder,
			second: eqOrder,
			want:   true,
		},
		{
			name:   "Orders with same location",
			first:  createValidOrder(kernel.RandomLocation(), 10),
			second: createValidOrder(kernel.RandomLocation(), 10),
			want:   false,
		},
		{
			name:   "Orders with different location",
			first:  createValidOrder(kernel.RandomLocation(), 10),
			second: createValidOrder(kernel.RandomLocation(), 10),
			want:   false,
		},
		{
			name:   "Orders with different volume",
			first:  createValidOrder(kernel.RandomLocation(), 5),
			second: createValidOrder(kernel.RandomLocation(), 10),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.first.Equals(tt.second))
			assert.Equal(t, tt.want, tt.second.Equals(tt.first))
			assert.True(t, tt.first.Equals(tt.first))
			assert.True(t, tt.second.Equals(tt.second))
		})
	}
}

func TestOrder_Assign(t *testing.T) {
	tests := []struct {
		name      string
		courierID uuid.UUID
		wantErr   error
	}{
		{
			name:      "Valid assign",
			courierID: uuid.New(),
		},
		{
			name:      "Nil courier id",
			courierID: uuid.Nil,
			wantErr:   errs.NewValueIsInvalidError("courierID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := createValidOrder(kernel.RandomLocation(), 10)

			got := o.Assign(tt.courierID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), got.Error())
				assert.Nil(t, o.CourierID())

				return
			}

			assert.Equal(t, &tt.courierID, o.CourierID())
			assert.Equal(t, order.StatusAssigned, o.Status())
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	o := createValidOrder(kernel.RandomLocation(), 10)
	err := o.Assign(uuid.New())
	require.Nil(t, err)

	err = o.Complete()
	require.Nil(t, err)
}

func TestOrder_CompleteAlreadyCompletedOrder(t *testing.T) {
	o := createValidOrder(kernel.RandomLocation(), 10)
	err := o.Assign(uuid.New())
	require.Nil(t, err)

	err = o.Complete()
	require.Nil(t, err)

	err = o.Complete()
	require.Error(t, err)
	require.ErrorIs(t, order.ErrOrderAlreadyCompleted, err)
}

func TestOrder_AssignAlreadyAssignedOrder(t *testing.T) {
	o := createValidOrder(kernel.RandomLocation(), 10)
	err := o.Assign(uuid.New())
	require.Nil(t, err)

	err = o.Assign(uuid.Nil)
	require.Error(t, err)
	require.ErrorIs(t, order.ErrCourierAlreadyAssigned, err)
}

func TestOrder_CompleteCreatedOrder(t *testing.T) {
	o := createValidOrder(kernel.RandomLocation(), 10)
	err := o.Complete()

	require.ErrorIs(t, order.ErrOrderNotAssigned, err)
}

func createValidOrder(location kernel.Location, volume int) *order.Order {
	o, err := order.NewOrder(uuid.New(), location, volume)

	if err != nil {
		panic(err)
	}

	return o
}
