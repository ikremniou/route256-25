package utils

import (
	"bufio"
	"strconv"
	"strings"
)

func ToNumberArray(rowString []string) ([]int, error) {
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

func ReadStringsAndConcat(reader *bufio.Reader, num int) string {
	var result = strings.Builder{}
	for i := 0; i < num; i++ {
		var str = ReadString(reader)
		result.Write([]byte(str))
	}

	return result.String()
}

func ReadNumber(reader *bufio.Reader) int {
	var valueStr, _ = reader.ReadString('\n')
	var value, _ = strconv.Atoi(strings.TrimSpace(valueStr))

	return value
}

func ReadString(reader *bufio.Reader) string {
	var valueStr, _ = reader.ReadString('\n')

	return strings.TrimRight(valueStr, "\n")
}

func ReadRowOfNumbers(reader *bufio.Reader) []int {
	var valueStr, _ = reader.ReadString('\n')
	var fields = strings.Fields(valueStr)
	var result, _ = ToNumberArray(fields)

	return result
}

func StrRowsToArrayOfNumbers(rows []string) [][]int {
	var result = make([][]int, len(rows))

	for i := 0; i < len(rows); i++ {
		var fields = strings.Fields(rows[i])
		result[i], _ = ToNumberArray(fields)
	}

	return result
}
