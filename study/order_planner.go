package study

import (
	"fmt"
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

func leftBinarySearch(plans []OrderPlannerCarSchedule, target int) int {
	var left = 0
	var right = len(plans) - 1
	var index = -1

	// 713837032 998778687 228866664 162098512 220833468 235613844 511379452 244654808 90465458 827889431
	// expected: []string{"3 -1 2 2 2 4 7 7 2 1"}
	// actual  : []string{"9 -1 2 2 2 4 3 7 2 9"}

	for left <= right {
		var middle = (left + right) / 2

		if plans[middle].Start <= target {
			index = middle
			left = middle + 1
		} else {
			right = middle - 1
		}
	}

	fmt.Println(index)

	if left >= len(plans) {
		return -1
	}

	var leftVal = plans[left]
	for left < len(plans) && plans[left] == leftVal {
		if plans[left].Start <= target && plans[left].End >= target && plans[left].Capacity > 0 {
			return left
		}

		left += 1
	}

	return -1
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

			var machineNumber = -1
			// binary search?
			// var index = leftBinarySearch(intervals, arrival)
			// if index >= 0 {
			// 	machineNumber = intervals[index].MachineNumber
			// 	intervals[index].Capacity--
			// }

			for k := 0; k < len(intervals) && machineNumber == -1; k++ {
				if intervals[k].Start <= arrival && intervals[k].End >= arrival && intervals[k].Capacity > 0 {
					machineNumber = intervals[k].MachineNumber
					intervals[k].Capacity--
				}
			}

			if j != 0 && j%100 == 0 {
				intervals = slices.DeleteFunc(intervals, func(interval OrderPlannerCarSchedule) bool {
					return interval.Capacity == 0 || interval.End < arrival
				})
			}

			currResult[arrivalsMap[arrival]] = strconv.Itoa(machineNumber)
		}

		result[i] = strings.Join(currResult, " ")
	}

	return result
}
