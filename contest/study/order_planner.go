package study

import (
	"slices"
	"strconv"
	"strings"
)

type OrderPlannerCarSchedule struct {
	Start         int
	End           int
	Capacity      int
	MachineNumber int
}

type OrderPlannerInput struct {
	NumberOfOrders int
	Arrival        []int
	NumberOfCars   int
	Schedule       []OrderPlannerCarSchedule
}

func OrderPlanner(input []OrderPlannerInput) []string {
	var result = make([]string, len(input))

	for i := 0; i < len(input); i++ {
		var intervals = input[i].Schedule
		var arrivalsMap = make(map[int]int)

		for j := 0; j < input[i].NumberOfOrders; j++ {
			arrivalsMap[input[i].Arrival[j]] = j
		}

		slices.Sort(input[i].Arrival)
		slices.SortFunc(intervals, func(a, b OrderPlannerCarSchedule) int {
			if a.Start < b.Start {
				return -1
			} else if a.Start > b.Start {
				return 1
			} else {
				return a.MachineNumber - b.MachineNumber
			}
		})

		var currResult = make([]string, input[i].NumberOfOrders)
		var machineCurrIndex = 0

		for j := 0; j < input[i].NumberOfOrders; j++ {
			var currentArrival = input[i].Arrival[j]

			var machineFoundIndex = -1
			for machineCurrIndex < input[i].NumberOfCars &&
				machineFoundIndex == -1 &&
				intervals[machineCurrIndex].Start <= currentArrival {

				var car = &intervals[machineCurrIndex]
				if car.Capacity > 0 && car.Start <= currentArrival && car.End >= currentArrival {
					car.Capacity -= 1
					machineFoundIndex = car.MachineNumber
				}

				if car.Capacity == 0 || car.End < currentArrival {
					machineCurrIndex += 1
				}
			}

			if machineFoundIndex == -1 {
				currResult[arrivalsMap[currentArrival]] = "-1"
			} else {
				currResult[arrivalsMap[currentArrival]] = strconv.Itoa(machineFoundIndex)
			}
		}

		result[i] = strings.Join(currResult, " ")
	}

	return result
}
