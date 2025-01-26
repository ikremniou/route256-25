package contest_test

import (
	"bufio"
	"ikremniou/route256/contest"
	"ikremniou/route256/study/utils"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func calculateSimilarStrings(reader *bufio.Reader) []int {
	var numberOfCases = utils.ReadNumber(reader)
	var result = make([]int, numberOfCases)

	for i := 0; i < numberOfCases; i++ {
		var numberOfStrings = utils.ReadNumber(reader)
		var myStrings = make([]string, numberOfStrings)

		for j := 0; j < numberOfStrings; j++ {
			myStrings[j] = strings.TrimSpace(utils.ReadString(reader))
		}

		result[i] = contest.CalculateSimilarity(myStrings)
	}

	return result
}

func TestSimilarStrings(t *testing.T) {
	const testName = "21"

	inputFile, err := os.Open(path.Join("test_data", "similar", testName))
	if err != nil {
		t.Fatal(err)
	}
	defer inputFile.Close()

	expectFile, err := os.ReadFile(path.Join("test_data", "similar", testName+".a"))
	if err != nil {
		t.Fatal(err)
	}

	var result = calculateSimilarStrings(bufio.NewReader(inputFile))

	var realExpert = make([]int, len(result))
	var expertReal = strings.Fields(string(expectFile))

	for i := 0; i < len(result); i++ {
		realExpert[i], _ = strconv.Atoi(expertReal[i])
	}

	assert.Equal(t, realExpert, result, "Input and Expect should be equal arrays")
}
