package study_test

import (
	"bufio"
	"encoding/json"
	"ikremniou/route256/study"
	"ikremniou/route256/study/utils"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readNumberOfVirusTestCases(reader *bufio.Reader) int {
	var numberOfCases = utils.ReadNumber(reader)

	return numberOfCases
}

func TestVirusFile(t *testing.T) {
	var testArray = []string{"1", "2", "3", "4", "5", "7", "8", "9", "10", "11"}

	for _, testName := range testArray {
		t.Run(testName, func(t *testing.T) {
			inputFile, err := os.Open(path.Join("test_data", "virus_files", testName))
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			expectFile, err := os.ReadFile(path.Join("test_data", "virus_files", testName+".a"))
			if err != nil {
				t.Fatal(err)
			}

			var reader = bufio.NewReader(inputFile)
			var numberOfInput = readNumberOfVirusTestCases(reader)
			var result = make([]int, numberOfInput)

			for i := 0; i < numberOfInput; i++ {
				var toRead = utils.ReadNumber(reader)
				var jsonString = utils.ReadStringsAndConcat(reader, toRead)
				var testCase = study.JsonFileSystem{}
				json.Unmarshal([]byte(jsonString), &testCase)

				result[i] = study.FindNumberOfVirusFiles(&testCase)
			}

			var expected, _ = utils.ToNumberArray(strings.Fields(string(expectFile)))

			assert.Equal(t, expected, result)
		})
	}
}
