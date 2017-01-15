package ride

import (
	"errors"
	"time"
)

var (
	ErrRiderAlreadyExist    = errors.New("Rider already exist")
	ErrRiderAlreadyDropped  = errors.New("Rider Already dropped")
	ErrActionNotPossible    = errors.New("Action not possible")
	ErrRiderNotFound        = errors.New("Rider with the ID not found")
	ErrRiderAlreadyPickedUp = errors.New("Rider already picked up")
)

func NewVehicle(capacity int64, location location) vehicle {
	return vehicle{
		Capacity:   capacity,
		Location:   location,
		Requestors: map[string]*requestor{},
		Riders:     map[string]*requestor{},
	}
}

func NewVehicleWithName(name string, capacity int64, location location) vehicle {
	return vehicle{
		ID:         name,
		Capacity:   capacity,
		Location:   location,
		Requestors: map[string]*requestor{},
		Riders:     map[string]*requestor{},
	}
}

type state struct {
	Active   bool
	Location location
}

type vehicle struct {
	State    state
	Capacity int64
	Location location
	ID       string
	// requestors always sorted by first one to drop.
	Requestors map[string]*requestor
	Riders     map[string]*requestor
	// req
	ExpectedLastDropTime time.Time
	// Path on which vehicle is currently travelling
	CurrentRoute routes
	// coputed values
	AllRoutes []routes
}

func (v *vehicle) addRequestor(r requestor) error {
	if _, ok := v.Riders[r.Identifier]; ok {
		return ErrRiderAlreadyExist
	}

	v.Riders[r.Identifier] = &r
	return nil
}

func (v vehicle) occupancyStatus() int64 {
	count := int64(0)
	for _, r := range v.Riders {
		if r.State == dropped {
			continue
		}
		count += r.Quantity
	}

	return count
}

func (v vehicle) remainingOccupancy() int64 {
	return v.Capacity - v.occupancyStatus()
}

func (v *vehicle) Drop(rider_id string) error {
	if v.Requestors == nil {
		v.Requestors = map[string]*requestor{}
	}

	found := false

	if rider, ok := v.Riders[rider_id]; ok {
		if rider.State == pickedUp {
			// set the state as dropped
			rider.State = dropped
			v.Requestors[rider_id] = rider

			delete(v.Riders, rider_id)
			found = true
		} else if rider.State == rideRequested {
			// remove the user from the pickup list
			// Cancel the ride.
			rider.State = requestCancelled
			v.Requestors[rider_id] = rider

			delete(v.Riders, rider_id)
			found = true
		} else if rider.State == dropped {
			// rider is removed if dropped, so this condition shouldnot happen
			return ErrRiderAlreadyDropped
		} else {
			return ErrActionNotPossible
		}
	}

	if found {
		// if the user is found and actioned, update the vehicle details
		return nil
	}

	return ErrRiderNotFound
}

func (v *vehicle) Pickup(rider_id string) error {
	if v.Requestors == nil {
		v.Requestors = map[string]*requestor{}
	}

	found := false

	if rider, ok := v.Riders[rider_id]; ok {
		if rider.State == pickedUp {
			return ErrRiderAlreadyPickedUp
		} else if rider.State == rideRequested {
			// remove the user from the pickup list
			// Cancel the ride.
			rider.State = pickedUp
			found = true
		} else if rider.State == dropped {
			// rider is removed if dropped, so this condition shouldnot happen
			return ErrRiderAlreadyDropped
		} else {
			return ErrActionNotPossible
		}
	}

	if found {
		return nil
	}

	return ErrRiderNotFound
}

func (v *vehicle) setStateForRider(rider_id string, stateRequested travelState) error {
	if rider, ok := v.Riders[rider_id]; ok {
		if rider.State == rideRequested {
			if stateRequested == pickedUp {
				rider.State = pickedUp
				return nil
			} else if stateRequested == rideRequested {
				return ErrActionNotPossible
			} else if stateRequested == dropped {
				rider.State = requestCancelled
				return nil
			}
			return nil
		} else if rider.State == pickedUp {
			return nil
		} else if rider.State == dropped {
			return nil
		} else if rider.State == requestCancelled {
			return nil
		} else {
			return ErrActionNotPossible
		}
	}

	return nil
}

func (v vehicle) GetRiderPins() pinList {
	pins := pinList{}
	for _, rider := range v.Riders {
		if rider.State == pickedUp {
			pins = append(pins, *NewPinFromRequestor(*rider, drop))
		} else if rider.State == rideRequested {
			pins = append(pins, *NewPinFromRequestor(*rider, pickup))
			pins = append(pins, *NewPinFromRequestor(*rider, drop))
		}

	}

	/*for _, r := range v.Requestors {
		pins = append(pins, *NewPinFromRequestor(*r, pickup))
		pins = append(pins, *NewPinFromRequestor(*r, drop))
	}*/

	return pins
}

func SegregateVehicles(vs []vehicle) ([]vehicle,[]vehicle) {
	empty := []vehicle{}
	occVehicle := []vehicle{}

	for _,v := range vs {
		if v.occupancyStatus() == 0 {
			empty = append(empty, v)
		} else {
			occVehicle = append(occVehicle, v)
		}
	}

	return empty, occVehicle
}