package study_test

import (
	"bufio"
	"ikremniou/route256/study"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readNumber(reader *bufio.Reader) int {
	var valueStr, _ = reader.ReadString('\n')
	var value, _ = strconv.Atoi(strings.TrimSpace(valueStr))

	return value
}

func readValues(reader *bufio.Reader) string {
	var valueStr, _ = reader.ReadString('\n')

	return strings.TrimRight(valueStr, "\n")
}

func getInputData(reader *bufio.Reader) []study.ResponseValidationInput {
	var inputSize = readNumber(reader)
	var input = make([]study.ResponseValidationInput, inputSize)

	for i := 0; i < inputSize; i++ {
		var numberOrValues = readNumber(reader)
		var actual = readValues(reader)
		var expected = readValues(reader)

		input[i] = study.ResponseValidationInput{
			NumberOfValues: numberOrValues,
			ActualStr:      actual,
			ExpectStr:      expected,
		}
	}
	return input
}

func TestResponseValidation(t *testing.T) {
	var testArray = []string{"1", "2", "3", "4", "5", "7", "8", "9",
		"10", "17", "28", "29", "31"}

	for _, testName := range testArray {
		t.Run(testName, func(t *testing.T) {

			inputFile, err := os.Open(path.Join("test_data", "validate_output", testName))
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()
			expectFile, err := os.ReadFile(path.Join("test_data", "validate_output", testName+".a"))
			if err != nil {
				t.Fatal(err)
			}

			var expectedResult = strings.Fields(string(expectFile))

			var input = getInputData(bufio.NewReader(inputFile))
			var result = study.ResponseValidation(input)

			assert.Equal(t, expectedResult, result, "Input and Expect should be equal 'yes/no' arrays")

		})
	}
}
