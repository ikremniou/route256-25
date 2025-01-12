package study_test

import (
	"bufio"
	"fmt"
	"ikremniou/route256/study"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getInputReadSize(reader *bufio.Reader) []string {
	var inputSize int
	fmt.Fscan(reader, &inputSize)

	return getInput(reader, inputSize)
}

func getInput(reader *bufio.Reader, inputSize int) []string {
	var input = make([]string, inputSize)

	for i := 0; i < inputSize; i++ {
		var number string
		fmt.Fscan(reader, &number)

		input[i] = number
	}

	return input
}

func TestRemoveDigit(t *testing.T) {
	for _, testName := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9",
		"10", "11", "12", "13", "16", "17", "18", "19", "20", "21", "22", "23", "24"} {
		t.Run(testName, func(t *testing.T) {
			inputFile, err := os.Open(path.Join("test_data", "remove_digit", testName))
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			expectFile, err := os.Open(path.Join("test_data", "remove_digit", testName+".a"))
			if err != nil {
				t.Fatal(err)
			}
			defer expectFile.Close()

			var inputReader = bufio.NewReader(inputFile)
			var input = getInputReadSize(inputReader)

			var expectReader = bufio.NewReader(expectFile)
			var expect = getInput(expectReader, len(input))

			var testResult = study.RemoveDigit(input)

			assert.Equal(t, expect, testResult, "Input and Expect should be equal arrays")
		})
	}
}
