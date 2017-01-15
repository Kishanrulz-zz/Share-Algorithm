package ride

func nextDrops(points pinList) []routes {
	// p := point{}
	// if rider is picked up, then drop points can be in any format.
	// handle pick up of non picked up riders first
	// if pick up done, find various permutation combination of drops.
	return []routes{}
}

func generateCombinations(pins pinList, length int) <-chan pinList {
	c := make(chan pinList)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan pinList) {
		defer close(c) // Once the iteration function is finished, we close the channel

		addLetter(c, NewPinList(), pins, length, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// AddLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
func addLetter(c chan pinList, combo pinList, pins pinList, length, max int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo pinList
	for _, ch := range pins {
		newCombo = combo.append(ch)
		if len(newCombo) == max && newCombo.valid() {
			c <- newCombo
		}
		addLetter(c, newCombo, pins.remove(pin(ch)), length-1, max)
	}
}

