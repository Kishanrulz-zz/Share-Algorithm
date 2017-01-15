package ride

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/skratchdot/open-golang/open"
	"time"
)

func AddVeh() (*vehicle, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	name := chomp(text)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle capacity: ")
	text, _ = reader.ReadString('\n')

	capacity, _ := strconv.ParseInt(chomp(text), 10, 64)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle current location: ")
	text, _ = reader.ReadString('\n')

	address := chomp(text)

	loc := &location{}
	var err error

	for {
		loc, err = NewLocationFromAddress(address)

		if err == nil {
			break
		}
		println(err.Error())
	}

	v := NewVehicleWithName(name, capacity, *loc)

	redisST.AddVehicle("blr", name, v.Location.Long, v.Location.Lat)

	_, err = redisST.InsertVehicles(v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func PickupRider() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	vehicle_name := chomp(text)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passenger name: ")
	text, _ = reader.ReadString('\n')
	rider_name := chomp(text)

	vehicles, err := redisST.FetchVehicleDetail(vehicle_name)
	if err != nil {
		return err
	}

	if len(vehicles) == 0 {
		return errors.New("Vehicle not found")
	}

	v := vehicles[0]
	err = v.Pickup(rider_name)
	if err != nil {
		return err
	}

	_, err = redisST.InsertVehicles(v)
	if err != nil {
		return err
	}

	return nil
}

func AddRequest() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Requestor name: ")
	text, _ := reader.ReadString('\n')
	name := chomp(text)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Requestor quantity: ")
	text, _ = reader.ReadString('\n')

	quantity, _ := strconv.ParseInt(chomp(text), 10, 64)

	reader = bufio.NewReader(os.Stdin)

	var pickLoc, dropLoc *location
	var err error
	for {
		fmt.Print("Enter Requestor pickup location: ")

		text, _ = reader.ReadString('\n')

		address := chomp(text)
		pickLoc, err = NewLocationFromAddress(address)
		if err == nil {
			break
		}
		println(err.Error())

	}
	for {

		fmt.Print("Enter Requestor drop location: ")
		text, _ = reader.ReadString('\n')

		address := chomp(text)

		dropLoc, err = NewLocationFromAddress(address)

		if err == nil {
			break
		}
		println(err.Error())
		return err
	}

	req := NewRequestor(name, quantity, *pickLoc, *dropLoc)

	vs, err := redisST.GetValidVehicleForRequestors(req)

	if err != nil {
		return err
	}

	now := time.Now()
	selRank, err := AssignVehicle(*req, vs)
	if err != nil {
		return err
	}

	pretty.Println(selRank.Route.toTimeString(now))
	_, err = redisST.InsertVehicles(selRank.V)
	if err != nil {
		return err
	}
	url, _ := selRank.Route.toMapAPI()

	open.Run(url)
	return nil
}

func RemoveVeh() (int64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	name := chomp(text)

	return redisST.RemoveVehicle("blr", name)
}

func chomp(command string) string {
	return strings.Trim(command, "\n")
}
