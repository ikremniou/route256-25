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
		for j := 0; j < input[i].NumberOfOrders; j++ {
			var arrival = input[i].Arrival[j]
			// 82520165 => to Machine 2 (but need to go to 9)
			// 61799209 => to Machine 9 (but need to go to 2)

			// binary search?
			var machineNumber = -1
			for k := 0; k < len(intervals) && machineNumber == -1; k++ {
				if intervals[k].Start <= arrival && intervals[k].End >= arrival && intervals[k].Capacity > 0 {
					machineNumber = intervals[k].MachineNumber
					intervals[k].Capacity--
				}
			}

			currResult[arrivalsMap[arrival]] = strconv.Itoa(machineNumber)
		}

		result[i] = strings.Join(currResult, " ")
	}

	return result
}
