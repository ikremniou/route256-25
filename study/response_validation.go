package study

import (
	"slices"
	"strconv"
	"strings"
)

type ResponseValidationInput struct {
	NumberOfValues int
	ActualStr      string
	ExpectStr      string
}

func toNumberArray(rowString []string) ([]int, error) {
	var numbers = make([]int, len(rowString))
	for i := 0; i < len(rowString); i++ {
		var number, err = strconv.Atoi(rowString[i])
		if err != nil {
			return nil, err
		}
		numbers[i] = number
	}

	return numbers, nil
}

func determineYesNo(actual string, expected string) string {
	if len(actual) != len(expected) {
		return "no"
	}

	actualNumbers, err := toNumberArray(strings.Split(actual, " "))
	if err != nil {
		return "no"
	}

	expectedNumbers, err := toNumberArray(strings.Split(expected, " "))
	if err != nil {
		return "no"
	}

	slices.Sort(actualNumbers)

	for i := 0; i < len(actualNumbers); i++ {
		if actualNumbers[i] != expectedNumbers[i] {
			return "no"
		}
	}

	return "yes"
}

func ResponseValidation(input []ResponseValidationInput) []string {
	var yesNo = make([]string, len(input))
	for i := 0; i < len(input); i++ {
		yesNo[i] = determineYesNo(input[i].ActualStr, input[i].ExpectStr)
	}

	return yesNo
}
