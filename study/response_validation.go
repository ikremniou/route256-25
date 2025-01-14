package study

import (
	"ikremniou/route256/study/utils"
	"slices"
	"strings"
)

type ResponseValidationInput struct {
	NumberOfValues int
	ActualStr      string
	ExpectStr      string
}

func determineYesNo(actual string, expected string) string {
	if len(actual) != len(expected) {
		return "no"
	}

	actualNumbers, err := utils.ToNumberArray(strings.Split(actual, " "))
	if err != nil {
		return "no"
	}

	expectedNumbers, err := utils.ToNumberArray(strings.Split(expected, " "))
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
