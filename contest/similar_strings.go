package contest

import "strings"

func buildStringOddEven(input string, isOdd bool) string {
	var builder = strings.Builder{}

	for i := 0; i < len(input); i++ {
		if isOdd && i%2 == 0 {
			builder.WriteRune(rune(input[i]))
		} else if !isOdd && i%2 == 1 {
			builder.WriteRune(rune(input[i]))
		}
	}

	return builder.String()
}

func CalculateSimilarity(input []string) int {
	var similarityMap = make(map[string]int)
	var result = 0

	for i := 0; i < len(input); i++ {
		var odd = buildStringOddEven(input[i], true)
		var even = buildStringOddEven(input[i], false)

		var oddKey = odd + "_odd"
		var evenKey = even + "_even"

		var oddValue, isOddFound = similarityMap[oddKey]
		var evenValue, isEvenFound = similarityMap[evenKey]

		if !isOddFound && odd != "" {
			similarityMap[oddKey] = 1
		} else if odd != "" {
			similarityMap[oddKey]++
		}

		if !isEvenFound && even != "" {
			similarityMap[evenKey] = 1
		} else if even != "" {
			similarityMap[evenKey]++
		}

		if oddValue > evenValue {
			result += oddValue
		} else {
			result += evenValue
		}
	}

	return result
}
