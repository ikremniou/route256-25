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

func getInputData(reader *bufio.Reader) []study.ResponseValidationInput {
	var inputSize = utils.ReadNumber(reader)
	var input = make([]study.ResponseValidationInput, inputSize)

	for i := 0; i < inputSize; i++ {
		var numberOrValues = utils.ReadNumber(reader)
		var actual = utils.ReadString(reader)
		var expected = utils.ReadString(reader)

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
