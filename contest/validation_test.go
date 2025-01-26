package contest_test

import (
	"bufio"
	"ikremniou/route256/contest"
	"ikremniou/route256/study/utils"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getProductValidationInput(reader *bufio.Reader) []contest.ProductValidationInput {
	var numberOfCases = utils.ReadNumber(reader)
	var result = make([]contest.ProductValidationInput, numberOfCases)

	for i := 0; i < numberOfCases; i++ {
		var numberOfProducts = utils.ReadNumber(reader)
		var products = make(map[string]string)
		var uniquePrices = make(map[string]bool)

		for i := 0; i < numberOfProducts; i++ {
			var productSplit = strings.Fields(utils.ReadString(reader))

			products[productSplit[0]] = productSplit[1]
			uniquePrices[productSplit[1]] = false
		}

		var theInputString = utils.ReadString(reader)

		result[i] = contest.ProductValidationInput{
			NumberOfProducts: numberOfProducts,
			Products:         products,
			ProductPrices:    uniquePrices,
			ValidationString: theInputString,
		}
	}

	return result
}

func TestValidationTask(t *testing.T) {
	const testName = "36"

	inputFile, err := os.Open(path.Join("test_data", "validation", testName))
	if err != nil {
		t.Fatal(err)
	}
	defer inputFile.Close()

	expectFile, err := os.ReadFile(path.Join("test_data", "validation", testName+".a"))
	if err != nil {
		t.Fatal(err)
	}

	var input = getProductValidationInput(bufio.NewReader(inputFile))
	var result = contest.ProductValidation(input)

	var expertReal = strings.Fields(string(expectFile))
	assert.Equal(t, expertReal, result, "Input and Expect should be equal arrays")
}
