package network

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const maxSlots = 16384

const stateFile = "/home/aneeshie/sandbox/slot.state"

var (
	mu        sync.Mutex
	usedSlots map[int]bool
)

func init() {
	usedSlots = load()
}

func load() map[int]bool {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return make(map[int]bool)
	}

	var slots []int
	if err := json.Unmarshal(data, &slots); err != nil {
		return make(map[int]bool)
	}

	m := make(map[int]bool, len(slots))
	for _, s := range slots {
		m[s] = true
	}
	return m
}

func save() error {
	slots := make([]int, 0, len(usedSlots))
	for s := range usedSlots {
		slots = append(slots, s)
	}

	data, err := json.Marshal(slots)
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0o644)
}

func NextAvailable() int {
	mu.Lock()
	defer mu.Unlock()

	for slot := range maxSlots {
		if !usedSlots[slot] {
			usedSlots[slot] = true
			if err := save(); err != nil {
				fmt.Println("slot_manager: failed to persist state:", err)
			}
			return slot
		}
	}

	return -1
}

func Release(slot int) {
	mu.Lock()
	defer mu.Unlock()

	delete(usedSlots, slot)
	if err := save(); err != nil {
		fmt.Println("slot_manager: failed to persist state:", err)
	}
}