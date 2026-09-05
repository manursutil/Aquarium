package main

import "aquarium/internal/simulation"

func fishByID(fish []simulation.FishSnapshot, id simulation.FishID) (simulation.FishSnapshot, bool) {
	for _, value := range fish {
		if value.ID == id {
			return value, true
		}
	}
	return simulation.FishSnapshot{}, false
}
