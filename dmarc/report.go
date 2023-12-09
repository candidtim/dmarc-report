package dmarc

import (
	"fmt"
	"sort"
	"time"
)

func DeduplicateFeedbacks(feedbacks []Feedback) []Feedback {
	uniqueFeedbacks := make([]Feedback, 0, len(feedbacks))
	seen := make(map[Feedback]struct{})

	for _, feedback := range feedbacks {
		if _, exists := seen[feedback]; !exists {
			seen[feedback] = struct{}{}
			uniqueFeedbacks = append(uniqueFeedbacks, feedback)
		}
	}

	return uniqueFeedbacks
}

func OrderByBeginTime(feedbacks []Feedback) {
	sort.Slice(feedbacks, func(i, j int) bool {
		return feedbacks[i].ReportMetadata.DateRange.Begin.Before(feedbacks[j].ReportMetadata.DateRange.Begin)
	})
}

func GroupByOrgName(feedbacks []Feedback) map[string][]Feedback {
	grouped := make(map[string][]Feedback)

	for _, feedback := range feedbacks {
		orgName := feedback.ReportMetadata.OrgName
		grouped[orgName] = append(grouped[orgName], feedback)
	}

	return grouped
}

func MergeAdjacentFeedbacks(feedbacks []Feedback) []Feedback {
	if len(feedbacks) <= 1 {
		return feedbacks
	}

	merged := make([]Feedback, 0, len(feedbacks))
	current := feedbacks[0]

	for i := 1; i < len(feedbacks); i++ {
		next := feedbacks[i]

		adjacent := current.ReportMetadata.DateRange.End.Add(2 * time.Second).After(next.ReportMetadata.DateRange.Begin)
		sameResult := current.Record.Row.PolicyEvaluated.DKIM == next.Record.Row.PolicyEvaluated.DKIM &&
			current.Record.Row.PolicyEvaluated.SPF == next.Record.Row.PolicyEvaluated.SPF

		if adjacent && sameResult {
			// Merge the current and next feedbacks
			current.ReportMetadata.DateRange.End = next.ReportMetadata.DateRange.End
		} else {
			// No merge, add the current to the result
			merged = append(merged, current)
			current = next
		}
	}

	// Add the last feedback to the result
	merged = append(merged, current)

	return merged
}

func FilterFeedbacks(feedbacks []Feedback, filter func(Feedback) bool) []Feedback {
	filtered := make([]Feedback, 0, len(feedbacks))

	for _, feedback := range feedbacks {
		if filter(feedback) {
			filtered = append(filtered, feedback)
		}
	}

	return filtered
}

func FormatStatus(status string, color bool) string {
	if status == "fail" {
		status = "FAIL"
	}

	if !color {
		return status
	}

	switch status {
	case "pass":
		return fmt.Sprintf("\x1b[32m%s\x1b[0m", status)
	case "FAIL":
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", status)
	default:
		return status
	}
}

func FormatFeedback(feedback Feedback, color bool) string {
	return fmt.Sprintf(
		"%s\t%s\t\t%s\t%s\t\t%s\t%s",
		feedback.ReportMetadata.DateRange.Begin.Format(time.DateTime),
		feedback.ReportMetadata.DateRange.End.Format(time.DateTime),
		FormatStatus(feedback.Record.Row.PolicyEvaluated.DKIM, color),
		FormatStatus(feedback.Record.Row.PolicyEvaluated.SPF, color),
		FormatStatus(feedback.Record.AuthResults.DKIM.Result, color),
		FormatStatus(feedback.Record.AuthResults.SPF.Result, color),
	)
}

func FormatHeader() string {
	return fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"Begin              ",
		"End                ",
		"Policy:",
		"DKIM",
		"SPF",
		"Auth:",
		"DKIM",
		"SPF",
	)
}
