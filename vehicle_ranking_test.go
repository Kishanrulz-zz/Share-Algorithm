package ride


import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"fmt"
	// "github.com/kr/pretty"
)

func TestVehicleRanking(t *testing.T) {
	assert := assert.New(t)

	// Rider2PickUP, _ := NewLocationFromAddress("41, Srinivagilu Main Rd, Ejipura, Bengaluru, Karnataka 560007")
	// Rider2Drop, _ := NewLocationFromAddress("IBM Ln, Embassy Golf Links Business Park, Challaghatta, Bengaluru, Karnataka 560071")

	car1CurrLoc, _ := NewLocationFromAddress("13.025190,77.636776")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop, errr := NewLocationFromAddress("13.004701,77.635832") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084
	fmt.Println("ERRRRRR::",errr)

	Rider2PickUP, _ := NewLocationFromAddress("13.010388,77.631283") //25, Sadashiva Temple Rd, KSFC Layout, Lingarajapuram, Bengaluru, Karnataka 560084
	Rider2Drop, _ := NewLocationFromAddress("13.001607,77.624073") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc, _ := NewLocationFromAddress("13.021925,77.634051")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop, _ := NewLocationFromAddress("13.011290,77.663083") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016

	vehicle1 := vehicle{
		ID: "khrm1",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*16),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*16),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	vehicle2 := vehicle{
		ID: "khrm2",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*21),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*21),
	}

	// req := requestor{
	// 	Identifier: "rider-2",
	// 	State: rideRequested,
	// 	Quantity: 1,
	// 	PickupLocation: *Rider2PickUP,
	// 	DropLocation: *Rider2Drop,
	// }

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop

	ranking := GetVehiclesRanking([]vehicle{vehicle1, vehicle2},req.Identifier, reqPickUpPin, reqDropPin)

	assert.Len(ranking, 2)
	assert.Equal("khrm1", ranking[0].V.ID)

	assert.True(true)
}