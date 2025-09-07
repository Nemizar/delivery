package courier_test

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCourier(t *testing.T) {
	type args struct {
		name     string
		speed    int
		location kernel.Location
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Valid courier",
			args: args{
				name:     "Courier",
				speed:    10,
				location: kernel.RandomLocation(),
			},
		},
		{
			name: "Invalid courier name",
			args: args{
				name:     "",
				speed:    10,
				location: kernel.RandomLocation(),
			},
			wantErr: errs.NewValueIsInvalidError("name"),
		},
		{
			name: "Invalid courier speed",
			args: args{
				name:     "Courier",
				speed:    0,
				location: kernel.RandomLocation(),
			},
			wantErr: errs.NewValueIsInvalidError("speed"),
		},
		{
			name: "Invalid location",
			args: args{
				name:     "Courier",
				speed:    10,
				location: kernel.Location{},
			},
			wantErr: errs.NewValueIsInvalidError("location"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := courier.NewCourier(tt.args.name, tt.args.speed, tt.args.location)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)

				return
			}

			assert.NotNil(t, got)
			assert.Equal(t, tt.args.name, got.Name())
			assert.Equal(t, tt.args.speed, got.Speed())
			assert.Equal(t, tt.args.location, got.Location())
			assert.Len(t, got.StoragePlaces(), 1)

			sp := got.StoragePlaces()[0]
			assert.Equal(t, "Bag", sp.Name())
			assert.Equal(t, 10, sp.TotalVolume())
		})
	}
}

func TestCourier_CanTakeOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderVolume int
		want        bool
	}{
		{
			name:        "Can",
			orderVolume: 2,
			want:        true,
		},
		{
			name:        "Can't",
			orderVolume: 20,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := order.NewOrder(uuid.New(), kernel.RandomLocation(), tt.orderVolume)
			require.Nil(t, err)

			c, err := courier.NewCourier("Courier", 2, kernel.RandomLocation())
			require.Nil(t, err)

			got, err := c.CanTakeOrder(o)
			require.Nil(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCourier_CanTakeOrderNil(t *testing.T) {
	c, err := courier.NewCourier("Courier", 2, kernel.RandomLocation())
	require.Nil(t, err)

	got, err := c.CanTakeOrder(nil)
	require.Error(t, err)
	assert.Equal(t, errs.NewValueIsInvalidError("order").Error(), err.Error())
	assert.False(t, got)
}

func TestCourier_TakeOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderVolume int
		wantErr     error
	}{
		{
			name:        "Valid order",
			orderVolume: 2,
		},
		{
			name:        "Big volume",
			orderVolume: 20,
			wantErr:     courier.ErrNoSuitablePlace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := order.NewOrder(uuid.New(), kernel.RandomLocation(), tt.orderVolume)
			require.Nil(t, err)

			c, err := courier.NewCourier("Courier", 2, kernel.RandomLocation())
			require.Nil(t, err)

			got := c.TakeOrder(o)
			if tt.wantErr != nil {
				assert.ErrorIs(t, got, tt.wantErr)

				return
			}

			require.Equal(t, tt.wantErr, err)
			assert.Nil(t, got)
		})
	}
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	tests := []struct {
		name            string
		courierLocation kernel.Location
		orderLocation   kernel.Location
		want            float64
		wantErr         error
	}{
		{
			name:            "Distance 1",
			courierLocation: createValidLocation(1, 1),
			orderLocation:   createValidLocation(2, 2),
			want:            1,
		},
		{
			name:            "Distance 3.5",
			courierLocation: createValidLocation(1, 1),
			orderLocation:   createValidLocation(4, 5),
			want:            3.5,
		},
		{
			name:            "Distance 9",
			courierLocation: createValidLocation(1, 1),
			orderLocation:   createValidLocation(10, 10),
			want:            9,
		},

		{
			name:            "Distance 0",
			courierLocation: createValidLocation(1, 1),
			orderLocation:   createValidLocation(1, 1),
			want:            0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := courier.NewCourier("Courier", 2, tt.courierLocation)
			require.Nil(t, err)

			got, err := c.CalculateTimeToLocation(tt.orderLocation)
			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func createValidLocation(x, y int) kernel.Location {
	l, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}

	return l
}
