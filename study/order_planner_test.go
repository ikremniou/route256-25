package study_test

import (
	"bufio"
	"ikremniou/route256/study"
	"ikremniou/route256/study/utils"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getOrderPlannerInputData(reader *bufio.Reader) []study.OrderPlannerInput {
	var inputSize = utils.ReadNumber(reader)
	var input = make([]study.OrderPlannerInput, inputSize)

	for i := 0; i < inputSize; i++ {
		var numberOfOrders = utils.ReadNumber(reader)
		var arrival = utils.ReadRowOfNumbers(reader)
		var numberOfCars = utils.ReadNumber(reader)
		var schedule = make([]study.OrderPlannerCarSchedule, numberOfCars)

		for j := 0; j < numberOfCars; j++ {
			var some = utils.ReadRowOfNumbers(reader)

			schedule[j] = study.OrderPlannerCarSchedule{
				Start:         some[0],
				End:           some[1],
				Capacity:      some[2],
				MachineNumber: j + 1,
			}
		}

		input[i] = study.OrderPlannerInput{
			NumberOfOrders: numberOfOrders,
			Arrival:        arrival,
			NumberOfCars:   numberOfCars,
			Schedule:       schedule,
		}
	}

	return input
}

func TestOrderPlanner(t *testing.T) {
	var testArray = []string{"1", "2", "3", "4", "5", "6"}

	for _, testName := range testArray {
		t.Run(testName, func(t *testing.T) {

			inputFile, err := os.Open(path.Join("test_data", "order_planner", testName))
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()
			expectFile, err := os.ReadFile(path.Join("test_data", "order_planner", testName+".a"))
			if err != nil {
				t.Fatal(err)
			}

			var expectTrimmed = strings.TrimSpace(string(expectFile))
			var realExpect = strings.Split(expectTrimmed, "\n")
			for i := 0; i < len(realExpect); i++ {
				realExpect[i] = strings.TrimSpace(realExpect[i])
			}

			var input = getOrderPlannerInputData(bufio.NewReader(inputFile))
			var result = study.OrderPlanner(input)

			assert.Equal(t, realExpect, result, "Input and Expect should be equal delivery array numbers")

		})
	}

}
