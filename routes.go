package ride

// Starting point is the car location
// first pin has the distance and time from the vehicle
// each pin next has the distance and time from the previous one
type routes struct {
	ID          string
	Pins        pinList
}
