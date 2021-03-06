package database

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/types"
)

func TestFindRanges(t *testing.T) {
	ranges := []types.RangeData{
		{
			StartLine:      0,
			StartCharacter: 3,
			EndLine:        0,
			EndCharacter:   5,
		},
		{
			StartLine:      1,
			StartCharacter: 3,
			EndLine:        1,
			EndCharacter:   5,
		},
		{
			StartLine:      2,
			StartCharacter: 3,
			EndLine:        2,
			EndCharacter:   5,
		},
		{
			StartLine:      3,
			StartCharacter: 3,
			EndLine:        3,
			EndCharacter:   5,
		},
		{
			StartLine:      4,
			StartCharacter: 3,
			EndLine:        4,
			EndCharacter:   5,
		},
	}

	m := map[types.ID]types.RangeData{}
	for i, r := range ranges {
		m[types.ID(i)] = r
	}

	for i, r := range ranges {
		actual := findRanges(m, i, 4)
		expected := []types.RangeData{r}
		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("unexpected findRanges result %d (-want +got):\n%s", i, diff)
		}
	}
}

func TestFindRangesOrder(t *testing.T) {
	ranges := []types.RangeData{
		{
			StartLine:      0,
			StartCharacter: 3,
			EndLine:        4,
			EndCharacter:   5,
		},
		{
			StartLine:      1,
			StartCharacter: 3,
			EndLine:        3,
			EndCharacter:   5,
		},
		{
			StartLine:      2,
			StartCharacter: 3,
			EndLine:        2,
			EndCharacter:   5,
		},
		{
			StartLine:      5,
			StartCharacter: 3,
			EndLine:        5,
			EndCharacter:   5,
		},
		{
			StartLine:      6,
			StartCharacter: 3,
			EndLine:        6,
			EndCharacter:   5,
		},
	}

	m := map[types.ID]types.RangeData{}
	for i, r := range ranges {
		m[types.ID(i)] = r
	}

	actual := findRanges(m, 2, 4)
	expected := []types.RangeData{ranges[2], ranges[1], ranges[0]}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("unexpected findRanges result (-want +got):\n%s", diff)
	}

}

func TestComparePosition(t *testing.T) {
	left := types.RangeData{
		StartLine:      5,
		StartCharacter: 11,
		EndLine:        5,
		EndCharacter:   13,
	}

	testCases := []struct {
		line      int
		character int
		expected  int
	}{
		{5, 11, 0},
		{5, 12, 0},
		{5, 13, 0},
		{4, 12, +1},
		{5, 10, +1},
		{5, 14, -1},
		{6, 12, -1},
	}

	for _, testCase := range testCases {
		if cmp := comparePosition(left, testCase.line, testCase.character); cmp != testCase.expected {
			t.Errorf("unexpected comparisonPosition result for %d:%d. want=%d have=%d", testCase.line, testCase.character, testCase.expected, cmp)
		}
	}
}
