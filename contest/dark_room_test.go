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

func getDarkRoomsInput(reader *bufio.Reader) []contest.DarkRoomInput {
	var numberOfRooms = utils.ReadNumber(reader)
	var result = make([]contest.DarkRoomInput, numberOfRooms)

	for i := 0; i < numberOfRooms; i++ {
		var some = utils.ReadRowOfNumbers(reader)
		result[i] = contest.DarkRoomInput{
			X: some[0],
			Y: some[1],
		}

	}

	return result
}

func TestDarkRoom(t *testing.T) {
	const testName = "1"

	inputFile, err := os.Open(path.Join("test_data", "dark_room", testName))
	if err != nil {
		t.Fatal(err)
	}
	defer inputFile.Close()

	// expectFile, err := os.ReadFile(path.Join("test_data", "dark_room", testName+".a"))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var input = getDarkRoomsInput(bufio.NewReader(inputFile))
	var result = contest.DarkRoom(input)

	var out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {
			fmt.Fprintln(out, result[i][j])
		}
	}
}
