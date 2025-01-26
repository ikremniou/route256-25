package contest_test

import (
	"bufio"
	"fmt"
	"ikremniou/route256/contest"
	"ikremniou/route256/study/utils"
	"os"
	"path"
	"testing"
)

func readBoxesInput(reader *bufio.Reader) []contest.BoxesInput {
	var numberOfCases = utils.ReadNumber(reader)
	var result = make([]contest.BoxesInput, numberOfCases)

	for i := 0; i < numberOfCases; i++ {
		var mAndN = utils.ReadRowOfNumbers(reader)
		var matrix = make([]string, mAndN[0])

		for j := 0; j < mAndN[0]; j++ {
			matrix[j] = utils.ReadString(reader)
		}

		result[i] = contest.BoxesInput{
			Matrix: matrix,
		}
	}

	return result
}

func TestBoxes(t *testing.T) {
	const testName = "1"

	inputFile, err := os.Open(path.Join("test_data", "boxes", testName))
	if err != nil {
		t.Fatal(err)
	}
	defer inputFile.Close()

	// expectFile, err := os.ReadFile(path.Join("test_data", "boxes", testName+".a"))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var input = readBoxesInput(bufio.NewReader(inputFile))

	var result = contest.Boxes(input)
	fmt.Println(result)

}
