package former

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type selection struct {
	SportEventID string
	MarketID     string
	OddID        string
	MarketType   int
	SportID      string
	TournamentID string
	IsLive       bool
}

// nolint:funlen // whyNoLint: its ok for many test cases
func TestExpressFormer(t *testing.T) {
	testCases := []struct {
		name                   string
		allElementsCount       int
		systemSize             []int
		expressesSeIDsExpected [][]string
	}{
		{
			name:             "all_4_system_size_3",
			allElementsCount: 4,
			systemSize:       []int{3},
			expressesSeIDsExpected: [][]string{
				{"1", "2", "3"},
				{"1", "2", "4"},
				{"1", "3", "4"},
				{"2", "3", "4"},
			},
		},
		{
			name:             "all_6_system_size_5",
			allElementsCount: 6,
			systemSize:       []int{5},
			expressesSeIDsExpected: [][]string{
				{"1", "2", "3", "4", "5"},
				{"1", "2", "3", "4", "6"},
				{"1", "2", "3", "5", "6"},
				{"1", "2", "4", "5", "6"},
				{"1", "3", "4", "5", "6"},
				{"2", "3", "4", "5", "6"},
			},
		},
		{
			name:             "all_6_system_size_4",
			allElementsCount: 6,
			systemSize:       []int{4},
			expressesSeIDsExpected: [][]string{
				{"1", "2", "3", "4"},
				{"1", "2", "3", "5"},
				{"1", "2", "3", "6"},
				{"1", "2", "4", "5"},
				{"1", "2", "4", "6"},
				{"1", "2", "5", "6"},
				{"1", "3", "4", "5"},
				{"1", "3", "4", "6"},
				{"1", "3", "5", "6"},
				{"1", "4", "5", "6"},
				{"2", "3", "4", "5"},
				{"2", "3", "4", "6"},
				{"2", "3", "5", "6"},
				{"2", "4", "5", "6"},
				{"3", "4", "5", "6"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			expresses := FormExpresses(generateSelections(testCase.allElementsCount), testCase.systemSize)
			assert.Equal(t, testCase.expressesSeIDsExpected, extractSeIDs(expresses))
		})
	}
}

// allElements = 36, systemSize  = 34
// BenchmarkExpressFormer-12    	    1842	    693259 ns/op	 2642186 B/op	     643 allocs/op
// allElements = 20, systemSize  = 8
// BenchmarkExpressFormer-12    	      22	  45630637 ns/op	149797019 B/op	  126002 allocs/op
// allElements = 5, systemSize  = 4
// BBenchmarkExpressFormer-12    	 1000000	      1080 ns/op	    2920 B/op	      11 allocs/op
func BenchmarkExpressFormer(b *testing.B) {
	const (
		allElements = 36
		systemSize  = 34
	)

	selections := generateSelections(allElements)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FormExpresses(selections, []int{systemSize})
	}
}

func TestIncrement(t *testing.T) {
	testCases := []struct {
		allElementsCount  int
		index             []int
		nextIndexExpected []int
		isIncremented     bool
	}{
		{allElementsCount: 6, index: []int{0, 1, 2, 3}, nextIndexExpected: []int{0, 1, 2, 4}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 1, 2, 4}, nextIndexExpected: []int{0, 1, 2, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 1, 2, 4}, nextIndexExpected: []int{0, 1, 2, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 1, 2, 5}, nextIndexExpected: []int{0, 1, 3, 4}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 1, 3, 4}, nextIndexExpected: []int{0, 1, 3, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 1, 4, 5}, nextIndexExpected: []int{0, 2, 3, 4}, isIncremented: true},
		{allElementsCount: 6, index: []int{0, 2, 3, 4}, nextIndexExpected: []int{0, 2, 3, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{1, 2, 3, 4}, nextIndexExpected: []int{1, 2, 3, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{1, 3, 4, 5}, nextIndexExpected: []int{2, 3, 4, 5}, isIncremented: true},
		{allElementsCount: 6, index: []int{2, 3, 4, 5}, nextIndexExpected: []int{2, 3, 4, 5}, isIncremented: false},
		{allElementsCount: 20, index: []int{2, 3, 4, 5}, nextIndexExpected: []int{2, 3, 4, 6}, isIncremented: true},
		{allElementsCount: 20, index: []int{0, 3, 17, 18}, nextIndexExpected: []int{0, 3, 17, 19}, isIncremented: true},
		{allElementsCount: 20, index: []int{0, 3, 17, 19}, nextIndexExpected: []int{0, 3, 18, 19}, isIncremented: true},
		{allElementsCount: 20, index: []int{0, 3, 18, 19}, nextIndexExpected: []int{0, 4, 5, 6}, isIncremented: true},
		{allElementsCount: 20, index: []int{16, 17, 18, 19}, nextIndexExpected: []int{16, 17, 18, 19}, isIncremented: false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%d", t.Name(), i), func(t *testing.T) {
			expressIterator := newExpressIterator(
				len(testCase.index),
				testCase.allElementsCount,
			)

			expressIterator.expressSelIndexes = testCase.index
			isIncremented := expressIterator.next()
			assert.Equal(t, testCase.isIncremented, isIncremented)
			assert.Equal(t, testCase.nextIndexExpected, expressIterator.expressSelIndexes)
		})
	}
}

func generateSelections(count int) []selection {
	selections := make([]selection, count)
	for i := 1; i <= count; i++ {
		selections[i-1] = selection{
			SportEventID: fmt.Sprintf("%d", i),
			MarketID:     "test_market_id",
			OddID:        "1",
			MarketType:   1,
			SportID:      "test_sport_id",
			TournamentID: "test_tournament_id",
			IsLive:       true,
		}
	}

	return selections
}

func extractSeIDs(expresses [][]selection) [][]string {
	result := make([][]string, len(expresses))

	for i, express := range expresses {
		ids := make([]string, len(express))
		for j, s := range express {
			ids[j] = s.SportEventID
		}

		result[i] = ids
	}

	return result
}
