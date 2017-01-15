package ride

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"

	"gopkg.in/redis.v5"
	"github.com/kr/pretty"
)

var (
	redisST                   = NewRedisStore("localhost:6379", "")
	REGION                    = "blr"
	ErrNoNearbyVehicle = errors.New("No near by vehicle found")
)

func Store() {

}

type distanceUnit string

const (
	KM distanceUnit = "km"
	M  distanceUnit = "m"
)

type store interface {
	FetchAllVehicles() []vehicle
	FetchAllByRadius(string, float32, distanceUnit) []vehicle
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string) *redisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return &redisStore{client}
}

func (r redisStore) AddVehicle(key, name string, long, lat float64) (int64, error) {
	geoLocation := &redis.GeoLocation{
		Name:      name,
		Longitude: long,
		Latitude:  lat,
	}

	intCmd := r.client.GeoAdd(key, geoLocation)
	return intCmd.Result()
}

func (r redisStore) FetchAllByRadius(key string, long, lat, radius float64, unit distanceUnit) ([]redis.GeoLocation, error) {
	radiusQuery := &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      string(unit),
		WithDist:  true,
		WithCoord: true,
		Sort:      "ASC",
	}
	geoLocations := r.client.GeoRadius(key, long, lat, radiusQuery)
	return geoLocations.Result()
}

func (r redisStore) FetchVehicleDetail(keys ...string) ([]vehicle, error) {
	results, err := r.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}
	vs := []vehicle{}
	for _, result := range results {
		v := vehicle{}
		res, _ := result.(string)
		vBuff := bytes.NewBufferString(res)
		dec := gob.NewDecoder(vBuff)
		err = dec.Decode(&v)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (r redisStore) InsertVehicles(vs ...vehicle) (string, error) {

	var vStr []interface{}
	for _, v := range vs {
		var vehicleBuff bytes.Buffer
		enc := gob.NewEncoder(&vehicleBuff)

		err := enc.Encode(v)
		if err != nil {
			fmt.Println("Err:::", err)
			return "", errors.New("Can't encode to gob")
		}
		vStr = append(vStr, v.ID, vehicleBuff.String())
	}

	return r.client.MSet(vStr...).Result()
}

func (r redisStore) PickupRider(vehicle_id, rider_id string) error {
	vehicles, err := redisST.FetchVehicleDetail(vehicle_id)
	if err != nil {
		return err
	}

	if len(vehicles) == 0 {
		return errors.New("Vehicle not found")
	}

	v := vehicles[0]
	err = v.Pickup(rider_id)
	if err != nil {
		return err
	}

	_, err = redisST.InsertVehicles(v)
	if err != nil {
		return err
	}

	return nil
}

func (r redisStore) RemoveVehicle(key, name string) (int64, error) {
	intCmd := r.client.ZRem(key, name)
	return intCmd.Result()
}

func (r redisStore) GetIDsByRadius(loc location) ([]string, error) {
	ids := []string{}
	locations, err := r.FetchAllByRadius(REGION, loc.Long, loc.Lat, RADIUS, KM)
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		ids = append(ids, location.Name)
	}
	return ids, nil
}

func (r redisStore) GetValidVehicleForRequestors(req *requestor) ([]vehicle, error) {
	ids, err := r.GetIDsByRadius(req.PickupLocation)
	if err != nil {
		log.Println("Err Assign vehicle", err)
		//return DeviationResult{}, err
		return nil, err
	}
	pretty.Println(len(ids),ids)
	vs, err := r.FetchVehicleDetail(ids...)
	if err != nil {
		log.Println("Err Assign vehicle", err)
		//return DeviationResult{}, err
		return nil, err
	}
	validV := []vehicle{}
	for _, v := range vs {
		println("v.Capacity - v.occupancyStatus()",v.Capacity - v.occupancyStatus())
		if req.Quantity <= (v.Capacity - v.occupancyStatus()) {
			validV = append(validV, v)
		}
	}
	if len(validV) == 0 {
		return nil, ErrNoNearbyVehicle
	}
	return validV, nil
}
