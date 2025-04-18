package chess

import (
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

func numNodesHelper(g *game, depth int) int {
	if depth <= 0 {
		return 1
	}

	mvs := g.legalMoves()

	i := 0

	for _, m := range mvs {
		gC := g.copy()
		gC.applyMoveUnchecked(m)
		i += numNodesHelper(&gC, depth-1)
	}

	return i
}

func numNodes(fen string, depth int) (int, error) {
	g, err := ParseFEN(fen)

	if err != nil {
		return 0, err
	}

	return numNodesHelper(g, depth), nil
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
	testCases := readTestCases()
	for i, tc := range testCases {
		val, err := numNodes(tc.Fen, tc.Depth)

		if err != nil {
			t.Errorf("Error occured during test case: %v", i+1)
		}

		if tc.Nodes != val {
			t.Errorf("NumNodes(%v, %v) = %v, Want: %v", tc.Fen, tc.Depth, val, tc.Nodes)
		}
	}
}
