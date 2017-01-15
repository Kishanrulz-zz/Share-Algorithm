package ride

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVehicleOccupancy(t *testing.T) {
	assert := assert.New(t)

	v := NewVehicle(4, location{Lat: 12.223234, Long: 77.23123})
	v.addRequestor(*NewRequestor("rider-1", 1, location{Lat: 12.223234, Long: 77.23123}, location{Lat: 12.223234, Long: 77.23123}))

	assert.Equal(int64(1), v.occupancyStatus())

	assert.True(true)
}

func TestVehicleRemainingOccupancy(t *testing.T) {
	assert := assert.New(t)

	v := NewVehicle(4, location{Lat: 12.223234, Long: 77.23123})
	v.addRequestor(*NewRequestor("rider-1", 1, location{Lat: 12.223234, Long: 77.23123}, location{Lat: 12.223234, Long: 77.23123}))

	assert.Equal(int64(3), v.remainingOccupancy())

	assert.True(true)
}

func TestDropRider(t *testing.T) {
	assert := assert.New(t)

	v := vehicle{
		Capacity: 4,
		Location: location{},
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
			},
			"rider-2": &requestor{
				Identifier: "rider-2",
				State:      pickedUp,
				Quantity:   1,
			},
		},
	}

	assert.Len(v.Riders, 2)
	assert.Len(v.Requestors, 0)

	err := v.Drop("rider-1")
	assert.NoError(err)

	assert.Len(v.Riders, 1)
	assert.Len(v.Requestors, 1)

	assert.True(true)
}

func TestPickupRider(t *testing.T) {
	assert := assert.New(t)

	v := vehicle{
		Capacity: 4,
		Location: location{},
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
			},
			"rider-2": &requestor{
				Identifier: "rider-2",
				State:      rideRequested,
				Quantity:   1,
			},
		},
	}

	assert.Len(v.Riders, 2)
	assert.Len(v.Requestors, 0)

	err := v.Pickup("rider-2")
	assert.NoError(err)

	assert.Len(v.Riders, 2)
	assert.Len(v.Requestors, 0)

	assert.Equal(pickedUp, v.Riders["rider-2"].State)

	assert.True(true)
}

func TestVehicleRiderStateChange(t *testing.T) {
	assert := assert.New(t)

	v := vehicle{
		Capacity: 4,
		Location: location{},
		Riders: map[string]*requestor{
			"rider-y": &requestor{
				Identifier: "rider-y",
				State:      pickedUp,
				Quantity:   1,
			},
			"rider-g": &requestor{
				Identifier: "rider-g",
				State:      rideRequested,
				Quantity:   1,
			},
		},
	}

	err := v.setStateForRider("rider-g", rideRequested)
	assert.Error(err)

	err = v.setStateForRider("rider-g", pickedUp)

	assert.NoError(err)
	assert.Equal(pickedUp, v.Riders["rider-g"].State)

	assert.True(true)
}

