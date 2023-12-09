package main

import (
	"fmt"
	"os"
	"path/filepath"

  "github.com/candidtim/dmarc-report/dmarc"
)

func listDir(directoryPath, pattern string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(directoryPath, pattern))
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		if fileInfo.Mode().IsRegular() {
			filePaths = append(filePaths, file)
		}
	}

	return filePaths, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide both a directory path and a domain name as arguments.")
		return
	}

	directoryPath := os.Args[1]
	domainName := os.Args[2]
	globPattern := fmt.Sprintf("*!%s!*", domainName)

	fileList, err := listDir(directoryPath, globPattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var feedbacks []dmarc.Feedback
	for _, filePath := range fileList {
		feedback, err := dmarc.ParseDMARCReport(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing file:", err)
		} else {
			feedbacks = append(feedbacks, feedback)
		}
	}

	fmt.Println("# DKIM report for domain:", domainName)

	feedbacks = dmarc.DeduplicateFeedbacks(feedbacks)
	dmarc.OrderByBeginTime(feedbacks)
	grouped := dmarc.GroupByOrgName(feedbacks)
	for orgName, feedbacks := range grouped {
		fmt.Println()
		fmt.Println("## Reporter:", orgName)
		fmt.Println()
		fmt.Println(dmarc.FormatHeader())
		merged := dmarc.MergeAdjacentFeedbacks(feedbacks)
		for _, feedback := range merged {
			fmt.Println(dmarc.FormatFeedback(feedback))
		}
	}
}
