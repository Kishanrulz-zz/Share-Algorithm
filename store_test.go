package ride

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/kr/pretty"
)

func TestAddAndFindByRadius(t *testing.T) {
	assert := assert.New(t)
	t_redisStore := NewRedisStore("localhost:6379", "")
	defer t_redisStore.client.Close()

	t_redisStore.AddVehicle("blr", "C1", 77.580643, 12.972442)
	t_result, err := t_redisStore.FetchAllByRadius("blr", 77.580643, 12.972442, 1, KM)

	assert.NoError(err)
	assert.Equal(1, len(t_result), "Count mismatch, not all cars fetched")

	// Update location of Car

	t_redisStore.AddVehicle("blr", "C1", 77.580643, 12.972442)
	t_result, err = t_redisStore.FetchAllByRadius("blr", 77.580643, 12.972442, 1, KM)

	assert.NoError(err)
	assert.Equal(1, len(t_result), "Count mismatch, not all cars fetched")

	t_redisStore.AddVehicle("blr", "C1", 27.580643, 42.972442)

	t_result, err = t_redisStore.FetchAllByRadius("blr", 77.580643, 12.972442, 1, KM)

	assert.NoError(err)
	assert.Equal(0, len(t_result), "Cars old location shouldnot exist")

	t_result, err = t_redisStore.FetchAllByRadius("blr", 27.580643, 42.972442, 1, KM)

	assert.NoError(err)
	assert.Equal(1, len(t_result), "Cars new locations should exist")

}

func TestInsertVehicle(t *testing.T) {
	assert := assert.New(t)

	v1 := "c1"

	v := vehicle{
		ID: v1,
		Capacity: 3,
		Location: location{
			Lat: 12.978273,		//lakshmipura bus stop
			Long: 77.631454,
		},
	}



	cmd := NewRedisStore("localhost:6379", "")

	t_vehicle := NewVehicle(4, location{"", 77.644396, 12.961543})

	cmd.AddVehicle("blr", v1, v.Location.Long, v.Location.Lat)

	t_vehicle.ID = "c2"
	s, e := cmd.InsertVehicles(v,t_vehicle)
	assert.Equal("OK", s)
	assert.NoError(e)

	vs, err := cmd.FetchVehicleDetail(v1,"c2")

	pretty.Println(vs)

	assert.Len(vs, 2)
	assert.Equal(v, vs[0])
	assert.Equal(t_vehicle, vs[1])
	assert.NoError(err)

}

func TestVehicleGetByIdsRadius(t *testing.T) {
	assert := assert.New(t)

	car1CurrLoc := NewLocationFromLatLong(13.025190,77.636776)	// CDG Platinum building, 5th Cross Rd, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043
	car2CurrLoc := NewLocationFromLatLong(13.043405,77.609656)      //Cafe Thulp, No.21/22, 2nd Cross Road, CPR Layout, Kammanahalli, Bengaluru, Karnataka 560084

	cmd := NewRedisStore("localhost:6379", "")

	cmd.AddVehicle("blr", "khrm1", car1CurrLoc.Long, car1CurrLoc.Lat)
	cmd.AddVehicle("blr", "khrm2", car2CurrLoc.Long, car2CurrLoc.Lat)
	reqLoc := NewLocationFromLatLong(13.025190,77.636776)
	vs, err := cmd.GetIDsByRadius(*reqLoc)
	assert.NoError(err)
	assert.NotEmpty(vs)
	assert.True(true)
}

func TestRedisStore_GetValidVehicleForRequestors(t *testing.T) {
	assert := assert.New(t)

	v_1 := "c1"
	v := vehicle{
		ID: v_1,
		Capacity: 3,
		Location: location{
			Lat: 12.978273,		//lakshmipura bus stop
			Long: 77.631454,
		},
	}



	cmd := NewRedisStore("localhost:6379", "")

	t_vehicle := vehicle{
		ID: "test",
		Capacity: 4,
		Location: location{
			Lat: 12.961543,		//lakshmipura bus stop
			Long: 77.644396,
		},
	}
	t_vehicle.ID = "c2"

	cmd.AddVehicle("blr", v_1,v.Location.Long, v.Location.Lat)
	cmd.AddVehicle("blr", t_vehicle.ID,t_vehicle.Location.Long, t_vehicle.Location.Lat)

	s, e := cmd.InsertVehicles(t_vehicle,v)
	assert.Equal("OK", s)
	assert.NoError(e)
	Rider3PickUP := NewLocationFromLatLong(12.938794, 77.629494, "close to shelton royale, koramangala")
	Rider3Drop := NewLocationFromLatLong(12.970949, 77.657897, "Suranjan Das Rd") //Prestige Milton Garden Apartment, Milton St, D Costa Layout, Cooke Town, Bengaluru, Karnataka 560005
	req := requestor{
		Identifier: "rider-3",
		State: rideRequested,
		Quantity: 4,
		PickupLocation: *Rider3PickUP,
		DropLocation: *Rider3Drop,
	}

	vs, err:= cmd.GetValidVehicleForRequestors(&req)
	assert.NoError(err)

	pretty.Println(len(vs),err,vs)

	assert.True(true)


}