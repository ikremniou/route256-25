package contest

import (
	"encoding/json"
	"strings"
)

type BoxesInput struct {
	Matrix []string
}

type BoxItem struct {
	Children []BoxItem
	Space    int
	Name     string
}

type Pair struct {
	X, Y int
}

func insideBox(i, j int, matrix []string, visited map[Pair]bool) *BoxItem {
	var realI = i + 1
	var realJ = j + 1

	if visited[Pair{i, j}] {
		return nil
	}

	visited[Pair{i, j}] = true
	var myNameBuilder = strings.Builder{}

	for matrix[realI][realJ] != '+' && matrix[realI][realJ] != '-' && matrix[realI][realJ] != '|' && matrix[realI][realJ] != '.' {
		myNameBuilder.WriteRune(rune(matrix[realI][realJ]))

		realJ++
	}

	var myName = myNameBuilder.String()

	var width = 0
	var height = 0

	realI = i + 1
	realJ = j + 1

	for matrix[realI][realJ] != '|' {
		width += 1
		realJ += 1
	}

	realJ = j + 1
	for matrix[realI][realJ] != '-' {
		height += 1
		realI += 1
	}

	var boxItem = &BoxItem{
		Children: make([]BoxItem, 0),
		Space:    height * width,
		Name:     myName,
	}

	for x := i + 1; x < i+1+height; x++ {
		for y := j + 1; y < j+1+width; y++ {
			if matrix[x][y] == '+' && matrix[x][y+1] == '-' && matrix[x+1][y] == '|' {
				var box = insideBox(x, y, matrix, visited)
				if box != nil {
					boxItem.Children = append(boxItem.Children, *box)
				}

			}
		}
	}

	return boxItem
}

func Boxes(input []BoxesInput) string {
	var result = make([]BoxItem, len(input))

	for inputIndex := 0; inputIndex < len(input); inputIndex++ {
		var thisBoxItem = BoxItem{
			Children: make([]BoxItem, 0),
			Space:    0,
			Name:     "",
		}
		var visited = make(map[Pair]bool)

		for i := 0; i < len(input[inputIndex].Matrix); i++ {
			var currLine = input[inputIndex].Matrix[i]
			var nextLine string = ""
			if i+1 < len(input[inputIndex].Matrix) {
				nextLine = input[inputIndex].Matrix[i+1]
			}

			for j := 0; j < len(currLine); j++ {
				if currLine[j] == '+' && len(currLine) > j+1 && currLine[j+1] == '-' && len(nextLine) > j && nextLine[j] == '|' {

					var insideBoxItem = insideBox(i, j, input[inputIndex].Matrix, visited)
					if insideBoxItem != nil {
						thisBoxItem.Children = append(thisBoxItem.Children, *insideBoxItem)
					}

				}
			}
		}

		result[inputIndex] = thisBoxItem
	}

	// to correct json!!!
	var res, _ = json.Marshal(result)
	return string(res)
}
