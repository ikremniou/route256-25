package contest

import (
	"fmt"
)

type DarkRoomInput struct {
	X int
	Y int
}

func DarkRoom(input []DarkRoomInput) [][]string {
	var result = make([][]string, len(input))

	for index, room := range input {
		var iteration = make([]string, 0)

		if room.X == 1 {
			iteration = append(iteration, "1")
			iteration = append(iteration, "1 1 R")
		} else if room.Y == 1 {
			iteration = append(iteration, "1")
			iteration = append(iteration, "1 1 D")
		} else if room.X >= room.Y {
			iteration = append(iteration, "2")
			iteration = append(iteration, "1 1 D")
			iteration = append(iteration, fmt.Sprintf("%d %d U", room.X, room.Y))
		} else {
			iteration = append(iteration, "2")
			iteration = append(iteration, "1 1 R")
			iteration = append(iteration, fmt.Sprintf("%d %d L", room.X, room.Y))
		}

		result[index] = iteration
	}

	return result
}
