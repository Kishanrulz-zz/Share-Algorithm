package ride

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/kr/pretty"
	"fmt"
	"github.com/skratchdot/open-golang/open"
)

func TestDistanceMatRouteKFC (t *testing.T) {
	assert := assert.New(t)

	carCurrLoc, _ := NewLocationFromAddress("Ganapathi Temple 18th Rd, 6th Block, Koramangala, Bengaluru, Karnataka 560030")
	Rider1Drop, _ := NewLocationFromAddress("2037, 1st Cross Rd, Kodihalli, Bengaluru, Karnataka 560008")
	Rider2PickUP, _ := NewLocationFromAddress("41, Srinivagilu Main Rd, Ejipura, Bengaluru, Karnataka 560007")
	Rider2Drop, _ := NewLocationFromAddress("IBM Ln, Embassy Golf Links Business Park, Challaghatta, Bengaluru, Karnataka 560071")
	Rider3Pickup, err := NewLocationFromAddress("1st Block Koramangala, Koramangala, Bengaluru, Karnataka 560034")
	Rider3Drop, _ := NewLocationFromAddress("Konen Agrahara, Konena Agrahara, Murgesh Pallya, Bengaluru, Karnataka 560017")

	assert.NoError(err)
	pretty.Println(Rider3Pickup)

	carCurrLoc2, _ := NewLocationFromAddress("No A, Floor,, 445, 18th Main Rd, 1st Stage, Koramangala, Bengaluru, Karnataka 560095")

	v := vehicle{
		Capacity: 4,
		Location: *carCurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute * 7),
				DirectDropTime:time.Now().Add(time.Minute * 16),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute * 25),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	reqPickUpPin := *NewPinFromRequestor(req, pickup)        // New pin for upcoming rider's pickup
	reqDropPin := *NewPinFromRequestor(req, drop)                // New pin for upcoming rider's drop

	out := make(chan DeviationResult)

	var directDropTime time.Time
	go func() {
		_, directDropTime, _ = calculateDeviation(v, req.Identifier, reqPickUpPin, reqDropPin, time.Now(), out)
	}()
	//pretty.Println("TEST PRINT::::", routesCalculated)
	d := <-out
	pretty.Println("Go-JEK:::::", d.Route)

	v = vehicle{
		Capacity: 4,
		Location: *carCurrLoc2,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute * 30),
				DirectDropTime:time.Now().Add(time.Minute * 25),

			},
		 "rider-2": &requestor{
			Identifier: "rider-2",
			State: rideRequested,
			Quantity: 1,
			PickupLocation: *Rider2PickUP,
			DropLocation: *Rider2Drop,
			DirectDropTime: directDropTime,
		},
	},
		ExpectedLastDropTime: d.ExpectedLastTime,
	}

	req = requestor{
		Identifier: "rider-3",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider3Pickup,
		DropLocation: *Rider3Drop,
	}

	reqPickUpPin =  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin =  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop
	go func() {
		calculateDeviation(v, req.Identifier,reqPickUpPin, reqDropPin, time.Now().Add(time.Minute * 2), out) // time.Now() * 4 is the time between cars current location and Rider-2's pickup
	}()
	route := <- out
	// pretty.Println("Out Channel", <-out)

	assert.Equal(" -> START () -> DROP (rider-1) -> PICK_UP (rider-2) -> DROP (rider-2) -> PICK_UP (rider-3) -> DROP (rider-3)", route.Route.toString())

	close(out)
	//pretty.Println("TEST PRINT::::", routesCalculated2)
}

func PendingDistanceMatRouteKFCMultipleVehicle(t *testing.T) {

	carCurrLoc, _ := NewLocationFromAddress("Ganapathi Temple 18th Rd, 6th Block, Koramangala, Bengaluru, Karnataka 560030")
	Rider1Drop, _ := NewLocationFromAddress("2037, 1st Cross Rd, Kodihalli, Bengaluru, Karnataka 560008")
	Rider2PickUP, _ := NewLocationFromAddress("41, Srinivagilu Main Rd, Ejipura, Bengaluru, Karnataka 560007")
	Rider2Drop, _ := NewLocationFromAddress("IBM Ln, Embassy Golf Links Business Park, Challaghatta, Bengaluru, Karnataka 560071")

	v := vehicle{
		Capacity: 4,
		Location: *carCurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*16),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*25),
	}

	// v2 := vehicle{
	// 	Capacity: 4,
	// 	Location: *carCurrLoc,
	// 	Riders: map[string]*requestor{
	// 		"rider-1": &requestor{
	// 			Identifier: "rider-1",
	// 			State:      pickedUp,
	// 			Quantity:   1,
	// 			DropLocation: *Rider1Drop,
	// 			PickupTime: time.Now().Add(-time.Minute*7),
	// 			DirectDropTime:time.Now().Add(time.Minute*16),

	// 		},
	// 	},
	// 	ExpectedLastDropTime: time.Now().Add(time.Minute*25),
	// }

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop

	out := make(chan DeviationResult)

	var directDropTime time.Time
	go func() {
		_, directDropTime, _ = calculateDeviation(v, req.Identifier,reqPickUpPin, reqDropPin, time.Now(), out)
	}()
	//pretty.Println("TEST PRINT::::", routesCalculated)
	d := <- out
	pretty.Println("Go-JEK:::::",d.Route)
}

func TestMultiVehicleA(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(13.025190,77.636776)	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(13.004701,77.635832) // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(13.010388,77.631283)
	Rider2Drop := NewLocationFromLatLong(13.001607,77.624073) //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(13.043405,77.609656)  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(13.011290,77.663083) // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016


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

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", "khrm1", car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", "khrm2", car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)

	vs := []vehicle{vehicle1, vehicle2}

	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)
	pretty.Println("Rank 0::", ranks[0].V.ID)

	pretty.Println("Ranks::", ranks)

	//devResult, err := AssignVehicles(req)
	assert.Equal(2,len(ranks))

	pretty.Println("devResult::::", ranks, "err:::??", err)

	assert.Equal("khrm1",ranks[0].V.ID)

	for i, rank := range ranks {
		path, _ := rank.Route.toMapAPI()
		pretty.Println("rank #",i,"ROUTE",rank.Route.toString())
		open.Run(path)
	}

}

func TestMultiVehicle1(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasure metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(12.959241, 77.654099, "Murgeshpalya corner") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(12.954519, 77.681743, "Hindustan aeronautics") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016


	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*9),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*9),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*20),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)

	vs := []vehicle{vehicle1, vehicle2}

	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)
	pretty.Println("Rank 0::", ranks[0].V.ID)

	pretty.Println("Ranks::", ranks)

	//devResult, err := AssignVehicles(req)
	assert.Equal(2,len(ranks))

	pretty.Println("devResult::::", ranks, "err:::??", err)

	assert.Equal("ibra",ranks[0].V.ID)

	for _, rank := range ranks {
		path, _ := rank.Route.toMapAPI()
		open.Run(path)
	}

}


//When the second rider and requestor have same pickup and drop location the vehicle is choosen.
func TestMultiVehicle2(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasure metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(12.959241, 77.654099, "Murgeshpalya corner") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(12.954519, 77.681743, "Hindustan aeronautics") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016

	Rider3PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider3Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005


	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*9),

			},
			"rider-2": &requestor{
				Identifier: "rider-2",
				State: rideRequested,
				Quantity: 1,
				PickupLocation: *Rider2PickUP,
				DropLocation: *Rider2Drop,
				PickupTime: time.Now().Add(-time.Minute*20),
				DirectDropTime:time.Now().Add(time.Minute*30),
			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*30),
	}

	req := requestor{
		Identifier: "rider-3",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider3PickUP,
		DropLocation: *Rider3Drop,
	}

	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*20),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)

	vs := []vehicle{vehicle1, vehicle2}

	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)
	/*pretty.Println("Rank 0::", ranks[0].V.ID)

	pretty.Println("Ranks::", ranks)*/

	//devResult, err := AssignVehicles(req)
	assert.Equal(2,len(ranks))

	//pretty.Println("devResult::::", ranks, "err:::??", err)

	assert.Equal("ibra",ranks[0].V.ID)

	for _, rank := range ranks {
		path, _ := rank.Route.toMapAPI()
		open.Run(path)
	}

}

func TestMultiVehicle7(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasur metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(12.959241, 77.654099, "Murgeshpalya corner") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(12.954519, 77.681743, "Hindustan aeronautics") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016

	Rider3PickUP := NewLocationFromLatLong(12.938794, 77.629494, "close to shelton royale, koramangala")
	Rider3Drop := NewLocationFromLatLong(12.970949, 77.657897, "Suranjan Das Rd") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005


	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*9),

			},
			"rider-2": &requestor{
				Identifier: "rider-2",
				State: rideRequested,
				Quantity: 1,
				PickupLocation: *Rider2PickUP,
				DropLocation: *Rider2Drop,
				PickupTime: time.Now().Add(-time.Minute*20),
				DirectDropTime:time.Now().Add(time.Minute*30),
			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*30),
	}

	req := requestor{
		Identifier: "rider-3",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider3PickUP,
		DropLocation: *Rider3Drop,
	}

	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*20),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)

	vs := []vehicle{vehicle1, vehicle2}

	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)
	pretty.Println("Rank 0::", ranks[0].V.ID)

	pretty.Println("Ranks::", ranks)

	//devResult, err := AssignVehicles(req)
	assert.Equal(2,len(ranks))

	pretty.Println("devResult::::", ranks, "err:::??", err)


	for _, rank := range ranks {
		fmt.Println("Time string",rank.Route.toTimeString(time.Now()) )
		path, _ := rank.Route.toMapAPI()
		open.Run(path)
	}
}

func TestMultiVehicle10(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasur metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(12.959241, 77.654099, "Murgeshpalya corner") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(12.954519, 77.681743, "Hindustan aeronautics") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016

	/*Rider3PickUP := NewLocationFromLatLong(12.938794, 77.629494, "close to shelton royale, koramangala")
	Rider3Drop := NewLocationFromLatLong(12.970949, 77.657897, "Suranjan Das Rd") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005*/



	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*9),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*30),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*20),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)

	vs := []vehicle{vehicle1, vehicle2}

	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)

	//devResult, err := AssignVehicles(req)
	assert.Equal(2,len(ranks))



	for i, rank := range ranks {
		fmt.Println("\n\n\nBEST Route for # ",i ,rank.Route.toTimeString(time.Now()) )
		path, _ := rank.Route.toMapAPI()
		open.Run(path)
	}

}

func TestAssignVehicle11(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasur metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043

	car1Rider1Drop := NewLocationFromLatLong(12.959241, 77.654099, "Murgeshpalya corner") // 6, Balakrishnappa Rd, Ramaswamipalya, Lingarajapuram, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084
	car2Rider1Drop := NewLocationFromLatLong(12.954519, 77.681743, "Hindustan aeronautics") // Service Rd, Govindpura, Dooravani Nagar, Bengaluru, Karnataka 560016

	/*Rider3PickUP := NewLocationFromLatLong(12.938794, 77.629494, "close to shelton royale, koramangala")
	Rider3Drop := NewLocationFromLatLong(12.970949, 77.657897, "Suranjan Das Rd") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005*/



	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car1Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*9),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*30),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *car2Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*20),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	vs := []vehicle{vehicle1, vehicle2}

	rank, err := AssignVehicle(req, vs)

	assert.NoError(err)

	fmt.Println("\n\n\nBEST",rank.Route.toTimeString(time.Now()) )
	path, _ := rank.Route.toMapAPI()
	open.Run(path)

}
func TestAssignVehicle13(t *testing.T) {

	assert := assert.New(t)

	//car1curr to car1rider1drop == 9min

	car1CurrLoc := NewLocationFromLatLong(12.975888, 77.626312, "Halasur metro")	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043



	car2CurrLoc := NewLocationFromLatLong(12.951583, 77.621532, "Near koramangala police station")  //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084

	Rider2PickUP := NewLocationFromLatLong(12.956931, 77.641527, "Opposite ramada encore domlur")
	Rider2Drop := NewLocationFromLatLong(12.980000, 77.656247, "Bagmane tech park") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005

	/*Rider3PickUP := NewLocationFromLatLong(12.938794, 77.629494, "close to shelton royale, koramangala")
	Rider3Drop := NewLocationFromLatLong(12.970949, 77.657897, "Suranjan Das Rd") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005*/



	vehicle1 := vehicle{
		ID: "ibra",
		Capacity: 4,
		Location: *car1CurrLoc,
		Riders: map[string]*requestor{},
		ExpectedLastDropTime: time.Now().Add(time.Minute*30),
	}


	vehicle2 := vehicle{
		ID: "pogba",
		Capacity: 4,
		Location: *car2CurrLoc,
		Riders: map[string]*requestor{},
		ExpectedLastDropTime: time.Now().Add(time.Minute*20),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", vehicle1.ID, car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", vehicle2.ID, car2CurrLoc.Long, car2CurrLoc.Lat)


	count , err := cmd.InsertVehicles(vehicle1, vehicle2)
	fmt.Println("Count", count, "Errr", err)

	vs := []vehicle{vehicle1, vehicle2}

	rank, err := AssignVehicle(req, vs)

	assert.NoError(err)

	fmt.Println("\n\n\nBEST",rank.Route.toTimeString(time.Now()) )
	path, _ := rank.Route.toMapAPI()
	open.Run(path)

}


