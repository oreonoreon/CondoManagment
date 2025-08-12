package services

import (
	"awesomeProject/internal/airbnbScraper"
	"awesomeProject/internal/entities"
	"awesomeProject/internal/excelCalendarScraper/models"
	"context"
	"errors"
	"github.com/johnbalvin/gobnb/search"
	"go.uber.org/zap"
	"time"
)

type ServiceBnB struct {
	storageAirbnbHost            StorageAirbnbHost
	storageAirbnbPrice           StorageAirbnbPrice
	storageAirbnbRoom            StorageAirbnbRoom
	storageCoordinatesAirbnbRoom StorageCoordinatesAirbnbRoom
	storageRatingAirbnbRoom      StorageRatingAirbnbRoom
}

type StorageAirbnbHost interface {
	CreateAirbnbHost(ctx context.Context, h entities.Host) (*entities.Host, error)
	GetAirbnbHost(ctx context.Context, id string) (*entities.Host, error)
}

type StorageAirbnbPrice interface {
	CreateAirbnbPrice(ctx context.Context, p entities.AirbnbPrice) (*entities.AirbnbPrice, error)
	GetAirbnbPrice(ctx context.Context, p entities.AirbnbPrice) (*entities.AirbnbPrice, error)
}

type StorageAirbnbRoom interface {
	CreateAirbnbRoom(ctx context.Context, r entities.AirbnbRoom) (*entities.AirbnbRoom, error)
	GetAirbnbRoom(ctx context.Context, roomID int64) (*entities.AirbnbRoom, error)
	UpdateAirbnbRoomUnderstandableType(ctx context.Context, roomID int64, understandableType string) (*entities.AirbnbRoom, error)
}

type StorageCoordinatesAirbnbRoom interface {
	CreateCoordinatesAirbnbRoom(ctx context.Context, c entities.CoordinatesAirbnbRoom) (*entities.CoordinatesAirbnbRoom, error)
	UpdateCoordinatesAirbnbRoom(ctx context.Context, roomID int64, locationName string) (*entities.CoordinatesAirbnbRoom, error)
	GetCoordinatesAirbnbRoom(ctx context.Context, roomID int64) (*entities.CoordinatesAirbnbRoom, error)
}

type StorageRatingAirbnbRoom interface {
	CreateRatingAirbnbRoom(ctx context.Context, r entities.RatingAirbnbRoom) (*entities.RatingAirbnbRoom, error)
}

func NewServiceBnB(storageAirbnbHost StorageAirbnbHost, storageAirbnbPrice StorageAirbnbPrice, storageAirbnbRoom StorageAirbnbRoom, storageCoordinatesAirbnbRoom StorageCoordinatesAirbnbRoom, storageRatingAirbnbRoom StorageRatingAirbnbRoom) *ServiceBnB {
	return &ServiceBnB{
		storageAirbnbHost:            storageAirbnbHost,
		storageAirbnbPrice:           storageAirbnbPrice,
		storageAirbnbRoom:            storageAirbnbRoom,
		storageCoordinatesAirbnbRoom: storageCoordinatesAirbnbRoom,
		storageRatingAirbnbRoom:      storageRatingAirbnbRoom,
	}
}

func (b *ServiceBnB) ScrapDataFromAirBnB(ctx context.Context, filter string, locationoption string, checkIn, checkOut string) ([]entities.AirbnbData, error) {
	in, err := models.TimeConvert(checkIn)
	if err != nil {
		return nil, err
	}
	out, err := models.TimeConvert(checkOut)
	if err != nil {
		return nil, err
	}

	data, err := airbnbScraper.ScrapDataFromAirBnB(locationoption, in, out)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return nil, err
	}

	//var filter = "studio"

	var allowed bool
	for _, v := range listOfAllowedUnderstandableType {
		if filter == v {
			allowed = true
			break
		}
	}
	if !allowed {
		zap.L().Debug("ScrapDataFromAirBnB filter is not set or wrong", zap.Error(errors.New("this filter is not allowed")))
	}

	airbnbData := make([]entities.AirbnbData, 0, 10)
	for _, v := range data {
		room, price, rating, coordinates, infoPrice := airbnbScraper.ModelConvertToEntities(v, in, out)

		ro, err := b.storageAirbnbRoom.GetAirbnbRoom(ctx, room.RoomID)
		if err != nil {
			zap.L().Error("ScrapDataFromAirBnB", zap.Error(err))
		}
		if ro != nil {
			room.UnderstandableType = ro.UnderstandableType
		} else {
			err = b.bnbDataWriteInDB(ctx, room, rating, coordinates, in, out)
			if err != nil {
				zap.L().Error("ScrapDataFromAirBnB", zap.Error(err))
			}
		}

		co, err := b.storageCoordinatesAirbnbRoom.GetCoordinatesAirbnbRoom(ctx, coordinates.RoomID)
		if err != nil {
			zap.L().Error("ScrapDataFromAirBnB", zap.Error(err))
		}
		if co != nil {
			coordinates = *co
		}

		//фильтрация
		if allowed {
			if room.UnderstandableType != filter {
				continue
			}
		}

		d := entities.NewAirbnbData(room, price, rating, coordinates, infoPrice)
		airbnbData = append(airbnbData, d)
	}
	return airbnbData, nil
}

func (b *ServiceBnB) BnbDataWriteInDB(ctx context.Context, data []search.Data, in, out time.Time) error {
	for _, v := range data {
		room, price, rating, coordinates, _ := airbnbScraper.ModelConvertToEntities(v, in, out)

		r, err := b.storageAirbnbRoom.GetAirbnbRoom(ctx, room.RoomID)
		if err != nil {
			zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
			return err
		}
		if r == nil {
			err = b.bnbDataWriteInDB(ctx, room, rating, coordinates, in, out)
			if err != nil {
				return err
			}
		}

		_, err = b.CreateAirbnbPrice(ctx, price)
		if err != nil {
			zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
			return err
		}
	}
	return nil
}

func (b *ServiceBnB) bnbDataWriteInDB(ctx context.Context, room entities.AirbnbRoom, rating entities.RatingAirbnbRoom, coordinates entities.CoordinatesAirbnbRoom, in, out time.Time) error {
	hostData, err := airbnbScraper.ScrapHostFromAirBnb(room.RoomID, in, out)
	if err != nil {
		return err
	}
	if hostData == nil {
		return errors.New("empty response")
	}

	host, err := b.storageAirbnbHost.GetAirbnbHost(ctx, hostData.Host.ID)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return err
	}
	if host == nil {
		h := airbnbScraper.ModelHostConvertToEntities(*hostData, false)

		host, err = b.storageAirbnbHost.CreateAirbnbHost(ctx, h)
		if err != nil {
			zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
			return err
		}

	}

	room.HostID = host.ID

	_, err = b.storageAirbnbRoom.CreateAirbnbRoom(ctx, room)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return err
	}

	_, err = b.storageRatingAirbnbRoom.CreateRatingAirbnbRoom(ctx, rating)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return err
	}
	_, err = b.storageCoordinatesAirbnbRoom.CreateCoordinatesAirbnbRoom(ctx, coordinates)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return err
	}
	return nil
}

func (b *ServiceBnB) CreateAirbnbPrice(ctx context.Context, price entities.AirbnbPrice) (*entities.AirbnbPrice, error) {
	p, err := b.storageAirbnbPrice.GetAirbnbPrice(ctx, price)
	if err != nil {
		zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
		return nil, err
	}
	if p == nil {
		p, err = b.storageAirbnbPrice.CreateAirbnbPrice(ctx, price)
		if err != nil {
			zap.L().Error("ServiceBnB.ScrapDataFromAirBnB", zap.Error(err))
			return nil, err
		}
	}
	return p, nil
}

var listOfAllowedLocationName []string = []string{"Title", "Halo"}
var listOfAllowedUnderstandableType []string = []string{"studio", "1-bedroom", "2-bedroom"}

func (b *ServiceBnB) UpdateLocationName(ctx context.Context, roomID int64, locationName string) (*entities.CoordinatesAirbnbRoom, error) {
	var allowed bool
	for _, v := range listOfAllowedLocationName {
		if locationName == v {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.New("this location is not allowed")
	}
	coordinate, err := b.storageCoordinatesAirbnbRoom.UpdateCoordinatesAirbnbRoom(ctx, roomID, locationName)
	if err != nil {
		return nil, err
	}
	return coordinate, nil
}

func (b *ServiceBnB) UpdateRoomUnderstandableType(ctx context.Context, roomID int64, understandableType string) (*entities.AirbnbRoom, error) {
	var allowed bool
	for _, v := range listOfAllowedUnderstandableType {
		if understandableType == v {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.New("this UnderstandableType is not allowed")
	}

	room, err := b.storageAirbnbRoom.UpdateAirbnbRoomUnderstandableType(ctx, roomID, understandableType)
	if err != nil {
		return nil, err
	}
	return room, nil
}
