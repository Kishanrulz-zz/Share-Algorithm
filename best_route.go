package ride

import (
	"fmt"
	"time"
	"math"
	"github.com/kr/pretty"
	"errors"
//	"github.com/skratchdot/open-golang/open"
)


var (
	SharingFactor = time.Duration(17)
	NormalFactor = time.Duration(10)
	ErrNoBestRoute = errors.New("No routes with max deviation of 1.7 found")
)
type DeviationResult struct {
	V                vehicle
	Route            pinList
	Deviation        time.Duration
	VehicleDeviation time.Duration
	ExpectedLastTime time.Time
	DirectDropTime   time.Time
	PickUpTime       time.Time
	DropTime         time.Time
	Error error
}

func calculateDeviation(v vehicle, reqID string,reqPickUpPin , reqDropPin pin, now time.Time, out chan DeviationResult) ( []pinList, time.Time, time.Duration ){

	riderPins := v.GetRiderPins()
	pins := pinList{}

	vehiclePin := NewPinFromVehicle(v)
	pins = riderPins.append(*vehiclePin)

	pins = append(pins, reqPickUpPin)
	pins = append(pins, reqDropPin)


	pinsWithMetrics, err := GDistanceMatrix(pins)
	if err != nil {
		out <- DeviationResult{
			V : v,
			Error: err,
		}
		return nil, now, 0
	}


	vehiclePinMatrix, _ := pinsWithMetrics.findByLatLongString(vehiclePin.latLongString())	// vehiclePinMatrix gives metrics of vehicle that is all the distances covered having vehicle the start point

	//now := time.Now()

	requestPickUpPinMatrix, _ := pinsWithMetrics.findByLatLongString(reqPickUpPin.latLongString())

	// from vehiclePinMatrix, calculating the direct time to pick up the upcoming rider
	reqBestPickUpTime := now.Add(vehiclePinMatrix.Distance[reqPickUpPin.latLongString()].Time)

	fmt.Println("reqPickUpTime:::", reqBestPickUpTime)

	// direct time from upcoming riders pickup to his drop
	reqBestDropTime := reqBestPickUpTime.Add(requestPickUpPinMatrix.Distance[reqDropPin.latLongString()].Time)

	// requestorTravel time if directly from pickup to drop
	//reqBestTravelTime := reqBestDropTime.Sub(reqPickUpTime)

	// adding reqDropTime to the upcoming rider's pickup pin
	reqPickUpPin.Rider.DirectDropTime = reqBestDropTime

	// adding reqDropTime to the upcoming rider's drop pin
	reqDropPin.Rider.DirectDropTime = reqBestDropTime

	riderPins = append(riderPins, reqPickUpPin, reqDropPin)

	routes_calculated := []pinList{}

	for combination := range generateCombinations(riderPins, riderPins.count()) {
		pretty.Println("Routes combination::", combination.toString())
		routes_calculated = append(routes_calculated, processEachPinWithMatrix(*vehiclePin, combination, pinsWithMetrics))
	}

	bestRouteDeviation := time.Duration(math.MaxInt64)
	bestRoute := pinList{}
	vehicleDeviation := time.Duration(math.MaxInt64)

	var stepTime,reqPickUpTime,reqDropTime time.Time

	for pID, pins := range routes_calculated {
		routeDeviation := time.Duration(0)
		stepTime = now
		for _, route := range pins {
			stepTime = stepTime.Add(route.TimeToCover)
			//fmt.Println("stepTime",stepTime, "route.State",route.NextState)
			if route.NextState == pickup {
				if route.Rider.Identifier == reqID {
					reqPickUpTime = stepTime
				}
			}
			if route.NextState == drop {
				if route.Rider.Identifier == reqID {
					reqDropTime = stepTime
				}

				// Deviation for driver who is going to be dropped at this step
				dev := stepTime.Sub(route.Rider.DirectDropTime)
				/*
				if dev < 0 {
					continue
				}
				*/
				routeDeviation += dev
				//fmt.Println("Rider", route.Rider.Identifier,"route.Rider.DirectDropTime",route.Rider.DirectDropTime,"dev",dev, "routeDeviation", routeDeviation)
				if route.Rider.DirectDropTime.Sub(now) * NormalFactor >  stepTime.Sub(now) * SharingFactor{
					break
				}
			}
		}

		pretty.Println("Route #", pID, "Vehicle", v.ID, "Route:::", pins.toTimeString(time.Now()))

		if routeDeviation < bestRouteDeviation {
			bestRouteDeviation = routeDeviation
			bestRoute = pins
			vehicleDeviation = stepTime.Sub(v.ExpectedLastDropTime)
			/*u, _ := pins.toMapAPI()
			open.Run(u)*/

			//fmt.Println("stepTime",stepTime,"deltaDeviation", vehicleDeviation,"reqBestDropTime",reqBestDropTime,"v.expectedLastDropTime",v.ExpectedLastDropTime,"bestRouteDeviation",bestRouteDeviation,"reqBestTravelTime",reqBestTravelTime)
		}
	}

	//pretty.Println("Best route index", bestRouteIndex, "Deviation::: ", bestRouteDeviation.Minutes(),"Delta Deviation:: ", vehicleDeviation.Minutes(), "BestRoute",bestRoute.toString(),"VEHICLEID::", v.ID)


	var errX error
	if bestRouteDeviation == time.Duration(math.MaxInt64) {
		errX = ErrNoBestRoute

	}


	out <- DeviationResult{
		v,
		bestRoute,
		bestRouteDeviation,
		vehicleDeviation,
		stepTime,
		reqBestDropTime,
		reqPickUpTime,
		reqDropTime,
		errX,
	}

	return routes_calculated, reqBestDropTime, vehicleDeviation
}
