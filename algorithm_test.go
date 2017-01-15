package ride

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePins(t *testing.T) {
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

	r1 := requestor{
		Identifier: "rider-3",
		State:      rideRequested,
		Quantity:   1,
	}

	pins := generatePins(v, r1)

	// _, err := pins.(pinList)

	x := makePinList(pins...)

	// assert.NoError(err)
	// assert.Len(pinList)

	assert.Equal(4, x.count())
	assert.Equal(3, x.nextStateCount(drop))
	assert.Equal(1, x.nextStateCount(pickup))

	assert.Equal(" -> DROP (rider-1) -> DROP (rider-2) -> PICK_UP (rider-3) -> DROP (rider-3)", x.toString(),
		"this is caused because maps doesnt guarantee order")

	assert.True(x.valid())

	wrongPins := []pin{pins[0], pins[1], pins[3], pins[2]}
	wrongListPin := makePinList(wrongPins...)

	assert.Equal(" -> DROP (rider-1) -> DROP (rider-2) -> DROP (rider-3) -> PICK_UP (rider-3)", wrongListPin.toString(),
		"this is caused because maps doesnt guarantee order")
	assert.False(wrongListPin.valid())
}

func TestPinListAppend(t *testing.T) {
	assert := assert.New(t)

	p1 := pin{NextState: drop, Rider: requestor{Identifier: "1", State: pickedUp, Quantity: 1}}
	p2 := pin{NextState: drop, Rider: requestor{Identifier: "2", State: pickedUp, Quantity: 1}}

	pinL := makePinList(p1, p2)

	p3 := pin{NextState: drop, Rider: requestor{Identifier: "3", State: pickedUp, Quantity: 1}}
	p := pinL.append(p3)

	assert.Equal(3, p.count())
}

func TestGenerateCombinations(t *testing.T) {
	p1 := pin{NextState: drop, Rider: requestor{Identifier: "1", State: pickedUp, Quantity: 1}}
	p2 := pin{NextState: drop, Rider: requestor{Identifier: "2", State: rideRequested, Quantity: 1}}
	p3 := pin{NextState: pickup, Rider: requestor{Identifier: "2", State: rideRequested, Quantity: 1}}

	pinL := makePinList(p1, p2, p3)

	for combination := range generateCombinations(pinL, pinL.count()) {
		fmt.Println(combination.toString()) // This is instead of process(combination)
	}
}

func TestPinListRemovePin(t *testing.T) {
	assert := assert.New(t)

	p1 := pin{NextState: drop, Rider: requestor{Identifier: "1", State: pickedUp, Quantity: 1}}
	p2 := pin{NextState: drop, Rider: requestor{Identifier: "2", State: rideRequested, Quantity: 1}}
	p3 := pin{NextState: pickup, Rider: requestor{Identifier: "2", State: rideRequested, Quantity: 1}}
	p4 := pin{NextState: pickup, Rider: requestor{Identifier: "4", State: rideRequested, Quantity: 1}}

	pinL := makePinList(p1, p2, p3, p4)
	_pinList := pinL.remove(p4)

	assert.Equal(3, _pinList.count())
}
