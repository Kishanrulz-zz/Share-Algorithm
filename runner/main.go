package main

import(
	"bufio"
	"fmt"
	"os"
	ride "bitbucket.org/z_team_gojek/ride-fair"
	"strings"
)

func chomp(command string) string {
	return strings.Trim(command, "\n")
}

func main() {
	var err error
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Choose from the following commands::\n ADD_VEH | RM_VEH | UPDT_VEH | REQ_RIDE | PICKUP | EXIT\n > ")
		text, _ := reader.ReadString('\n')

		err = nil

		switch(chomp(text)) {
		case "ADD_VEH", "1":
			_, err = ride.AddVeh()
		case "RM_VEH", "2":
			_, err = ride.RemoveVeh()
		case "UPDT_VEH", "3":
		case "REQ_RIDE", "4":
			err = ride.AddRequest()
		case "CNCL_RIDE", "":
		case "PICKUP", "5":
			// "vehicle_name", "rider_id"
			err = ride.PickupRider()
			fmt.Println("ERR: ", err.Error())
		case "EXIT", "7":
			os.Exit(0)
		default:
			fmt.Println(text)
		}
		if err == nil {
			fmt.Println("Successfully executed")
		} else {
			fmt.Println("Oops!! an error occured ", err.Error())
		}
    }
}