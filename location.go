package ride

import (
	"errors"
	"googlemaps.github.io/maps"
	"context"
)

const (
	//apiKeyOld = "AIzaSyCyIeVJaMJIL72mmuH8w20HNpypp-klD_4"
	//apiKey = "AIzaSyBlZ2uYfYw7t44rK1eby50C7G6OLbURAd8"
	//apiKey = "AIzaSyB8yoTSFMGF4fr_-K_QOrKR5ZpUSbwRzTY"
	apiKey = "AIzaSyB4Wy09j3GjgkE0bqsV-h4eI0sw2JPfapI"
)

func NewLocationFromLatLong(lat, long float64, address ...string) *location {
	return &location{Long: long, Lat: lat}
}

func NewLocationFromAddress(address string) (*location, error) {
	var client *maps.Client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))

	resp, err := client.Geocode(context.Background(), &maps.GeocodingRequest{
		Address: address,
	})

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("Address not found")
	}

	return NewLocationFromLatLong(resp[0].Geometry.Location.Lat, resp[0].Geometry.Location.Lng), nil
}

type location struct {
	Address string  `json:"address"`
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
}