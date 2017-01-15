package ride

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)


func TestAddress(t *testing.T) {
	//car1CurrLoc, err := NewLocationFromAddress("Kamanahalli Main Road, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043")
	////////car1CurrLoc, err := NewLocationFromAddress("Jalvayu Vihar Water Tank, Jal Vayu Vihar, Bengaluru, Karnataka 560043")

	assert := assert.New(t)

	car1CurrLoc := NewLocationFromLatLong(12.994540, 77.684735, "Kamanahalli Main Road, HRBR Layout 3rd Block, HRBR Layout, Kalyan Nagar, Bengaluru, Karnataka 560043")

	//open.Run("https://www.google.co.in/maps/dir/Ganapathi+Temple+18th+Rd,+6th+Block,+Koramangala,+Bengaluru,+Karnataka+560030/Team+Royal's+Apartment,+Venkat+Reddy+Layout,+6th+Block,+Koramangala,+Bengaluru,+Karnataka+560047/41,+Srinivagilu+Main+Rd,+Ejipura,+Bengaluru,+Karnataka+560007")

	fmt.Println(car1CurrLoc)
	//assert.NoError(err)
	assert.NotNil(car1CurrLoc)
}