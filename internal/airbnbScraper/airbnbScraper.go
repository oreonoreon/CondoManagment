package airbnbScraper

import (
	"awesomeProject/internal/entities"
	"errors"
	"github.com/johnbalvin/gobnb"
	"github.com/johnbalvin/gobnb/details"
	"github.com/johnbalvin/gobnb/search"
	"sort"
	"time"
)

func ModelConvertToEntities(data search.Data, checkIn time.Time, checkOut time.Time) (entities.AirbnbRoom, entities.AirbnbPrice, entities.RatingAirbnbRoom, entities.CoordinatesAirbnbRoom, entities.AdditionalInfoBnbPrice) {
	images := make([]string, 0, 3)

	for _, v := range data.Images {
		images = append(images, v.URL)
	}
	var price float32
	if data.Price.Unit.Discount == 0 {
		price = data.Price.Unit.Amount
	} else {
		price = data.Price.Unit.Discount
	}

	return entities.AirbnbRoom{
			RoomID:             data.RoomID,
			Badges:             data.Badges,
			Name:               data.Name,
			Title:              data.Title,
			TypeOfEstate:       data.Type,
			Kind:               data.Kind,
			Category:           data.Category,
			Images:             images,
			UnderstandableType: "",
			HostID:             "",
		},
		entities.AirbnbPrice{
			RoomID:       data.RoomID,
			Price:        price,
			CheckIn:      checkIn,
			CheckOut:     checkOut,
			ScrapingDate: time.Now(),
		},
		entities.RatingAirbnbRoom{
			RoomID:      data.RoomID,
			Value:       data.Rating.Value,
			ReviewCount: data.Rating.ReviewCount,
		},
		entities.CoordinatesAirbnbRoom{
			RoomID:       data.RoomID,
			Latitude:     data.Coordinates.Latitude,
			Longitude:    data.Coordinates.Longitud,
			LocationName: "",
		},
		entities.AdditionalInfoBnbPrice{
			Qualifier:      data.Price.Unit.Qualifier,
			CurrencySymbol: data.Price.Unit.CurrencySymbol,
		}

}

func ModelHostConvertToEntities(d details.Data, trustable bool) entities.Host {
	return entities.Host{
		ID:        d.Host.ID,
		Name:      d.Host.Name,
		Trustable: trustable,
	}
}

var titleNayang = search.CoordinatesInput{
	Ne: search.CoordinatesValues{
		Latitude: 8.086705448481512,
		Longitud: 98.30054396206339,
	},
	Sw: search.CoordinatesValues{
		Latitude: 8.083702740588947,
		Longitud: 98.2988994418385,
	},
}

var halo = search.CoordinatesInput{
	Ne: search.CoordinatesValues{
		Latitude: 8.085326302000913,
		Longitud: 98.30236034421057,
	},
	Sw: search.CoordinatesValues{
		Latitude: 8.081843061342493,
		Longitud: 98.29885218696285,
	},
}

func ScrapDataFromAirBnB(location string, in, out time.Time) ([]search.Data, error) {

	client := gobnb.NewClient("THB", nil)
	// zoom value from 1 - 20, so from the "square" like I said on the coorinates
	// This represents how much zoom there is on this square.
	zoomvalue := 18
	check := search.Check{
		In:  in,
		Out: out,
	}

	var coords search.CoordinatesInput

	switch location {
	case "Halo":
		coords = halo
	case "Title":
		coords = titleNayang
	default:
		return nil, errors.New("location is not available")
	}

	results, err := client.SearchAll(zoomvalue, coords, check)
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		vi := results[i].Price.Unit.Amount
		vj := results[j].Price.Unit.Amount
		if results[i].Price.Unit.Discount != 0 {
			vi = results[i].Price.Unit.Discount
		}
		if results[j].Price.Unit.Discount != 0 {
			vj = results[j].Price.Unit.Discount
		}
		return vi < vj
	})

	//clean result from empty structs
	cleanResults := make([]search.Data, 0, 10)
	for _, result := range results {
		if result.RoomID != 0 {
			cleanResults = append(cleanResults, result)
		}
	}

	return cleanResults, nil
}

func ScrapHostFromAirBnb(roomID int64, in, out time.Time) (*details.Data, error) {
	client := gobnb.NewClient("THB", nil)

	check := search.Check{
		In:  in,
		Out: out,
	}

	results, err := client.DetailsFromRoomID(roomID, check)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
