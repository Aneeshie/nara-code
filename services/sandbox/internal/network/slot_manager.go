package network

var usedSlots map[int]bool

const maxSlots = 16384

func init() {
	usedSlots = make(map[int]bool)
}

func NextAvailable() int {
	for slot := range maxSlots {
		if !usedSlots[slot] {
			usedSlots[slot] = true // mark as claimed
			return slot
		}
	}

	return -1
}

func Release(slot int){
	delete(usedSlots, slot)
}


