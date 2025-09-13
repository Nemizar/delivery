package courierrepo

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
)

func DomainToDTO(courier *courier.Courier) CourierDTO {
	dto := CourierDTO{
		ID:    courier.Id(),
		Name:  courier.Name(),
		Speed: courier.Speed(),
		Location: LocationDTO{
			X: courier.Location().X(),
			Y: courier.Location().Y(),
		},
	}

	sp := make([]*StoragePlaceDTO, len(courier.StoragePlaces()))
	for i, storagePlace := range courier.StoragePlaces() {
		sp[i] = &StoragePlaceDTO{
			ID:          storagePlace.Id(),
			Name:        storagePlace.Name(),
			TotalVolume: storagePlace.TotalVolume(),
			OrderID:     storagePlace.OrderID(),
			CourierID:   courier.Id(),
		}
	}

	dto.StoragePlaces = sp

	return dto
}

func DTOToDomain(dto CourierDTO) *courier.Courier {
	sp := make([]*courier.StoragePlace, len(dto.StoragePlaces))
	for i, storagePlaceDTO := range dto.StoragePlaces {
		sp[i] = courier.RestoreStoragePlace(storagePlaceDTO.ID, storagePlaceDTO.Name, storagePlaceDTO.TotalVolume, storagePlaceDTO.OrderID)
	}

	l, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)

	return courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, l, sp)
}
