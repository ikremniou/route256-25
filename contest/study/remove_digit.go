package study

func RemoveDigit(input []string) []string {
	var result = make([]string, len(input))

	for i := 0; i < len(input); i++ {
		if len(input[i]) == 1 {
			result[i] = "0"
			continue
		}

		var previousDigit = 10
		var indexOfMinDigit = 0
		for pos, digitRune := range input[i] {
			var realDigit = int(digitRune - '0')
			if previousDigit < realDigit {
				indexOfMinDigit = pos
				break
			}

			previousDigit = realDigit
		}

		if indexOfMinDigit > 0 {
			result[i] = input[i][:indexOfMinDigit-1] + input[i][indexOfMinDigit:]
		} else {
			result[i] = input[i][:len(input[i])-1]
		}
	}

	return result
}
