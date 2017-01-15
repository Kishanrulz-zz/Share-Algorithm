package ride

import(
	"time"
)

type travelState string

const (
	pickedUp         travelState = "PICKED_UP"
	rideRequested    travelState = "RIDE_REQUESTED"
	dropped          travelState = "DROPPED"
	requestCancelled travelState = "REQUEST_CANCELLED"
)

func NewRequestor(id string, quantity int64, pickupLocation, dropLocation location) *requestor {
	return &requestor{
		Identifier:     id,
		State:          rideRequested,
		Quantity:       quantity,
		PickupLocation: pickupLocation,
		DropLocation:   dropLocation,
		RequestTime:    time.Now(),
	}
}

type requestor struct {
	Identifier       string
	State            travelState
	Quantity         int64
	PickupTime       time.Time
	DropTime         time.Time
	PickupLocation   location
	DropLocation     location
	ExpectedPickUpTime time.Time
	ExpectedDropTime time.Time
	RequestTime      time.Time
				       // seconds
	ActualTravelTime int64 // directTime
	WaitTime         time.Duration
	TotalTime        time.Duration

	DirectDropTime   time.Time
	ElapsedTime      time.Duration
	RerouteTime      time.Time
	DeviationTime    time.Duration
}