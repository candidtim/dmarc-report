package dmarc

import (
	"reflect"
	"testing"
	"time"
)

//
// DeduplicateFeedbacks
//

func TestDeduplicateFeedbacks(t *testing.T) {
	// empty input
	runDeduplicateTest(t, []Feedback{}, []Feedback{})

	// two unique items
	uniqueItems := []Feedback{
		{ReportMetadata: ReportMetadata{OrgName: "example.com", ReportID: "123"}},
		{ReportMetadata: ReportMetadata{OrgName: "example.com", ReportID: "456"}},
	}
	runDeduplicateTest(t, uniqueItems, uniqueItems)

	// two identical items
	inputDuplicates := []Feedback{
		{ReportMetadata: ReportMetadata{OrgName: "example.com", ReportID: "123"}},
		{ReportMetadata: ReportMetadata{OrgName: "example.com", ReportID: "123"}},
	}
	expectedDeduplicated := []Feedback{
		{ReportMetadata: ReportMetadata{OrgName: "example.com", ReportID: "123"}},
	}
	runDeduplicateTest(t, inputDuplicates, expectedDeduplicated)
}

func runDeduplicateTest(t *testing.T, input, expected []Feedback) {
	result := DeduplicateFeedbacks(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

//
// OrderByBeginTime
//

func TestOrderByBeginTime(t *testing.T) {
	unorderedFeedbacks := []Feedback{
		{ReportMetadata: ReportMetadata{DateRange: DateRange{Begin: time.Date(2023, time.January, 3, 0, 0, 0, 0, time.Local)}}},
		{ReportMetadata: ReportMetadata{DateRange: DateRange{Begin: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)}}},
		{ReportMetadata: ReportMetadata{DateRange: DateRange{Begin: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local)}}},
	}

	OrderByBeginTime(unorderedFeedbacks)

	expectedOrder := []time.Time{
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
		time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local),
		time.Date(2023, time.January, 3, 0, 0, 0, 0, time.Local),
	}

	actualOrder := make([]time.Time, len(unorderedFeedbacks))
	for i, feedback := range unorderedFeedbacks {
		actualOrder[i] = feedback.ReportMetadata.DateRange.Begin
	}

	if !reflect.DeepEqual(actualOrder, expectedOrder) {
		t.Errorf("Expected order %v, got %v", expectedOrder, actualOrder)
	}
}

//
// GroupByOrgName
//

func TestGroupByOrgName(t *testing.T) {
	feedbacks := []Feedback{
		{ReportMetadata: ReportMetadata{OrgName: "GroupA"}},
		{ReportMetadata: ReportMetadata{OrgName: "GroupB"}},
		{ReportMetadata: ReportMetadata{OrgName: "GroupA"}},
	}

	grouped := GroupByOrgName(feedbacks)

	// Expected result
	expected := map[string][]Feedback{
		"GroupA": {
			{ReportMetadata: ReportMetadata{OrgName: "GroupA"}},
			{ReportMetadata: ReportMetadata{OrgName: "GroupA"}},
		},
		"GroupB": {
			{ReportMetadata: ReportMetadata{OrgName: "GroupB"}},
		},
	}

	if !reflect.DeepEqual(grouped, expected) {
		t.Errorf("Expected result %v, got %v", expected, grouped)
	}
}

//
// MergeAdjacentFeedbacks
//

func newAdjFeedback(begin, end time.Time) Feedback {
	return Feedback{
		ReportMetadata: ReportMetadata{
			DateRange: DateRange{
				Begin: begin,
				End:   end,
			},
		},
	}
}

func TestMergeAdjacentFeedbacks_OneFeedback(t *testing.T) {
	feedbacks := []Feedback{
		newAdjFeedback(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local)),
	}

	result := MergeAdjacentFeedbacks(feedbacks)

	expected := feedbacks

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result %v, got %v", expected, result)
	}
}

func TestMergeAdjacentFeedbacks_TwoAdjacent(t *testing.T) {
	feedbacks := []Feedback{
		newAdjFeedback(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local)),
		newAdjFeedback(time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 3, 0, 0, 0, 0, time.Local)),
		newAdjFeedback(time.Date(2023, time.January, 5, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 6, 0, 0, 0, 0, time.Local)),
	}

	result := MergeAdjacentFeedbacks(feedbacks)

	expected := []Feedback{
		newAdjFeedback(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 3, 0, 0, 0, 0, time.Local)),
		newAdjFeedback(time.Date(2023, time.January, 5, 0, 0, 0, 0, time.Local), time.Date(2023, time.January, 6, 0, 0, 0, 0, time.Local)),
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result %v, got %v", expected, result)
	}
}
