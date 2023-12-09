package main

import (
	"flag"
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

var printUsage = func() {
	fmt.Printf("Usage: %s [OPTIONS] DIRECTORY DOMAIN\n\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	var failuresOnly bool
	var colorOutput bool

	flag.Usage = printUsage
	flag.BoolVar(&failuresOnly, "f", false, "Show only the failures")
	flag.BoolVar(&colorOutput, "c", true, "Show color output")

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		printUsage()
		return
	}

	directoryPath := args[0]
	domainName := args[1]
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

	fmt.Println("DMARC report for domain:", domainName)

	feedbacks = dmarc.DeduplicateFeedbacks(feedbacks)
	dmarc.OrderByBeginTime(feedbacks)
	grouped := dmarc.GroupByOrgName(feedbacks)

	for orgName, feedbacks := range grouped {
		fmt.Println("\nReporter:", orgName)

		if failuresOnly {
			feedbacks = dmarc.FilterFeedbacks(feedbacks, dmarc.Feedback.HasFailures)
		}

		if len(feedbacks) == 0 {
			fmt.Println("pass")
			continue
		}

		fmt.Println(dmarc.FormatHeader())
		merged := dmarc.MergeAdjacentFeedbacks(feedbacks)
		for _, feedback := range merged {
			if failuresOnly && !feedback.HasFailures() {
				continue
			} else {
				fmt.Println(dmarc.FormatFeedback(feedback, colorOutput))
			}
		}
	}
}
