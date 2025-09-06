package courier_test

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoragePlace_New(t *testing.T) {
	type args struct {
		name        string
		totalVolume int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Valid storage values",
			args: args{
				name:        "Backpack",
				totalVolume: 15,
			},
		},
		{
			name: "Name with spaces",
			args: args{
				name:        "   ",
				totalVolume: 10,
			},
			wantErr: errs.NewValueIsInvalidError("name"),
		},
		{
			name: "Empty name",
			args: args{
				name:        "",
				totalVolume: 10,
			},
			wantErr: errs.NewValueIsInvalidError("name"),
		},
		{
			name: "Zero volume",
			args: args{
				name:        "Backpack",
				totalVolume: 0,
			},
			wantErr: errs.NewValueIsInvalidError("totalVolume"),
		},
		{
			name: "Volume less than zero",
			args: args{
				name:        "Backpack",
				totalVolume: -1,
			},
			wantErr: errs.NewValueIsInvalidError("totalVolume"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := courier.NewStoragePlace(tt.args.name, tt.args.totalVolume)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)

				return
			}

			a := tt.args
			assert.NotEqual(t, uuid.Nil, got.Id())
			assert.Equal(t, a.totalVolume, got.TotalVolume())
			assert.Equal(t, a.name, got.Name())
			assert.Nil(t, got.OrderID())
		})
	}
}

func TestStoragePlace_Equals(t *testing.T) {
	s1, err := courier.NewStoragePlace("Backpack", 10)
	require.Nil(t, err)

	s2, err := courier.NewStoragePlace("Backpack", 10)
	require.Nil(t, err)

	assert.False(t, s1.Equals(s2))
	assert.False(t, s1.Equals(nil))
	assert.True(t, s1.Equals(s1))

	assert.False(t, s2.Equals(s1))
	assert.False(t, s2.Equals(nil))
	assert.True(t, s2.Equals(s2))
}

func TestStoragePlace_Store(t *testing.T) {
	type args struct {
		orderID uuid.UUID
		volume  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Valid order",
			args: args{
				orderID: uuid.New(),
				volume:  10,
			},
		},
		{
			name: "Nil order id",
			args: args{
				orderID: uuid.Nil,
				volume:  10,
			},
			wantErr: errs.NewValueIsInvalidError("orderID"),
		},
		{
			name: "Zero volume",
			args: args{
				orderID: uuid.New(),
				volume:  0,
			},
			wantErr: errs.NewValueIsInvalidError("volume"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := courier.NewStoragePlace("Backpack", 15)
			require.Nil(t, err)

			got := storage.Store(tt.args.orderID, tt.args.volume)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), got.Error())
				fmt.Println(storage.OrderID())
				assert.Nil(t, storage.OrderID())

				return
			}

			assert.Nil(t, got)
			assert.Equal(t, &tt.args.orderID, storage.OrderID())
		})
	}
}

func TestStoragePlace_Store_IsOccupied(t *testing.T) {
	storage, err := courier.NewStoragePlace("Backpack", 15)
	require.Nil(t, err)

	got := storage.Store(uuid.New(), 10)
	assert.Nil(t, got)

	got = storage.Store(uuid.New(), 10)
	assert.Equal(t, courier.ErrCanNotStoreOrder.Error(), got.Error())
}

func TestStoragePlace_CanStore(t *testing.T) {
	type args struct {
		name        string
		totalVolume int
	}
	tests := []struct {
		name    string
		volume  int
		args    args
		want    bool
		wantErr error
	}{
		{
			name:   "Same storage and order volume",
			volume: 10,
			args: args{
				name:        "Backpack",
				totalVolume: 10,
			},
			want: true,
		},
		{
			name:   "Storage more order volume",
			volume: 10,
			args: args{
				name:        "Backpack",
				totalVolume: 15,
			},
			want: true,
		},
		{
			name:   "Order more storage volume",
			volume: 15,
			args: args{
				name:        "Backpack",
				totalVolume: 10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := courier.NewStoragePlace(tt.args.name, tt.args.totalVolume)
			require.Nil(t, err)

			got, err := s.CanStore(tt.volume)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.False(t, got)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStoragePlace_Clear_Occupied(t *testing.T) {
	s, err := courier.NewStoragePlace("Backpack", 15)
	require.Nil(t, err)

	err = s.Store(uuid.New(), 10)
	require.Nil(t, err)

	err = s.Clear()
	require.Nil(t, err)

	require.Nil(t, s.OrderID())
}

func TestStoragePlace_Clear_Empty(t *testing.T) {
	s, err := courier.NewStoragePlace("Backpack", 15)
	require.Nil(t, err)

	err = s.Clear()
	require.Error(t, err)
	require.Equal(t, courier.ErrStoreAlreadyClear.Error(), err.Error())

	require.Nil(t, s.OrderID())
}
