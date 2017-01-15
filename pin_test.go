package ride

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPinLatLongString(t *testing.T) {
	assert := assert.New(t)

	t_requestor := *NewRequestor("rider-1", 1,
		*NewLocationFromLatLong(12.975928, 77.638986),
		*NewLocationFromLatLong(12.959969, 77.641068))

	t_pin := *NewPinFromRequestor(t_requestor, drop)
	t_pinList := makePinList(t_pin)

	assert.Equal(t_pin.latLongString(), "12.959969,77.641068")

	assert.Equal(t_pinList.latLongList(), []string{"12.959969,77.641068"})

	assert.True(true)
}
