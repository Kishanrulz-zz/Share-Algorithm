package ride

import (
	"context"
	//"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"gopkg.in/redis.v5"
	"log"
)

type redisGeo []redis.GeoLocation

func GDistanceMatrix( pins pinList) (pinList, error) {
	disReq := &maps.DistanceMatrixRequest{}

	destination := pins.latLongList()

	disReq.Origins = destination
	disReq.Destinations = destination
	disReq.Mode = "driving"
	disReq.Units = "metric"
	disReq.DepartureTime = "now"
	//calling maps API
	var client *maps.Client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Println("Error GDistanceMatrix", err)
		return nil, err
	}

	resp, err := client.DistanceMatrix(context.Background(), disReq)
	if err != nil {
		log.Println("Error GDistanceMatrix", err)
		return nil, err
	}
	// pretty.Println("Google response:::::",resp)



	latLongAddressMap := map[string]string{}
	for i, address := range resp.OriginAddresses {
		pins[i].Location.Address = address
		latLongAddressMap[pins[i].latLongString()] = address
	}

	pinsWithMetrix := []pin{}
	for i, _p := range pins {
		_p.Distance = tdMetrixToLatLongMap{}
		for _j, _r := range resp.Rows[i].Elements {
			_p.Distance[pins[_j].latLongString()] = tdMetrix{
				Distance: _r.Distance.Meters,
				Time:     _r.DurationInTraffic,
			}
		}
		pinsWithMetrix = append(pinsWithMetrix, _p)
	}

	//pretty.Println("PIN METRICS ::: ",pinsWithMetrix)

	return pinsWithMetrix, nil

}

func processEachPinWithMatrix(vehiclePin pin, riderPins pinList, matrixPins pinList) pinList {
	pins := addVehiclePinWithRider(vehiclePin, riderPins)

	pinsWithDistance := NewPinList()

	prevPin := pins[0]

	for _, p_pins := range pins {
		prev_pins_matrix, err := matrixPins.findByLatLongString(prevPin.latLongString())
		if err != nil {
			continue
		}
		p_pins.MetersToCover = int64(prev_pins_matrix.Distance[p_pins.latLongString()].Distance)
		p_pins.TimeToCover = prev_pins_matrix.Distance[p_pins.latLongString()].Time
		address, _ := matrixPins.findByLatLongString(p_pins.latLongString())
		p_pins.Location.Address = address.Location.Address

		pinsWithDistance = append(pinsWithDistance, p_pins)

		prevPin = p_pins
	}

	return pinsWithDistance
}

func addVehiclePinWithRider(vehiclePin pin, riderPins pinList) pinList {
	pins := NewPinList()
	pins = append(pins, vehiclePin)
	pins = append(pins, riderPins...)
	return pins
}

