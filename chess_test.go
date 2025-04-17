package chess 

import (
	"chess/game"
	"encoding/json"
	"io"
	"os"
	"testing"
)

type testCase struct {
	Depth int    `json:"depth"`
	Nodes int    `json:"nodes"`
	Fen   string `json:"fen"`
}

func readTestCases() []testCase {
	file, err := os.Open("perft_positions.json")
	defer file.Close()

	if err != nil {
		panic(err.Error())
	}

	jsonVals, err := io.ReadAll(file)

	if err != nil {
		panic(err.Error())
	}

	testCases := make([]testCase, 0)

	err = json.Unmarshal(jsonVals, &testCases)

	if err != nil {
		panic(err.Error())
	}

	return testCases
}

func TestPerftTests(t *testing.T) {

	// name := "Gladys"
	// want := regexp.MustCompile(`\b` + name + `\b`)
	// msg, err := Hello("Gladys")
	// if !want.MatchString(msg) || err != nil {
	// 	t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	// }

	testCases := readTestCases()
	for i, tc := range testCases {
		val, err := game.NumNodes(tc.Fen, tc.Depth)

		if err != nil {
			t.Errorf("Error occured during test case: %v", i+1)
		}

		if tc.Nodes != val {
			t.Errorf("NumNodes(%v, %v) = %v, Want: %v", tc.Fen, tc.Depth,val,  tc.Nodes)
		}
	}
}
