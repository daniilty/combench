package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	argsLen  = 3
	exitCode = 1
)

var (
	reNum          *regexp.Regexp
	hundredPercent = 100.0
)

func initRegexp() error {
	var err error

	reNum, err = regexp.Compile(`[0-9|0-9\.0-9]+`)
	if err != nil {
		return err
	}

	return nil
}

func checkArgs() error {
	if len(os.Args) != argsLen {
		return fmt.Errorf("Usage: benchcompare [old.txt] [new.txt]")
	}

	return nil
}

func readFilesFromArgs() ([]string, []string, error) {
	oldFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		return []string{}, []string{}, err
	}

	newFile, err := os.ReadFile(os.Args[2])
	if err != nil {
		return []string{}, []string{}, err
	}

	splittedOldFile := strings.Split(string(oldFile), "\n")
	splittedNewFile := strings.Split(string(newFile), "\n")

	return splittedOldFile, splittedNewFile, nil
}

func getBenchlines() (string, string, error) {
	oldFileLines, newFileLines, err := readFilesFromArgs()
	if err != nil {
		return "", "", err
	}

	oldBenchLine := getBenchline(oldFileLines)
	if oldBenchLine == "" {
		return "", "", fmt.Errorf("old file does not contain benchmark related line")
	}

	newBenchLine := getBenchline(newFileLines)
	if newBenchLine == "" {
		return "", "", fmt.Errorf("new file does not contain benchmark related line")
	}

	return oldBenchLine, newBenchLine, nil
}

func getBenchline(fileLines []string) string {
	for _, l := range fileLines {
		if !strings.HasPrefix(l, "Benchmark") {
			continue
		}

		return l
	}

	return ""
}

func compareResults(old string, new string) error {
	oldSlice := reNum.FindAllString(old, -1)
	if len(oldSlice) == 0 {
		return fmt.Errorf("no numeric data inside old benchmark results")
	}

	newSlice := reNum.FindAllString(new, -1)
	if len(newSlice) == 0 {
		return fmt.Errorf("no numeric data inside new benchmark results")
	}

	parsedOldSlice, parsedNewSlice, err := getParsedSlices(oldSlice, newSlice)
	if err != nil {
		return err
	}

	formattedTOpsDiff := getParsedDiff(parsedNewSlice, parsedOldSlice, 1)
	formattedNsPerOpsDiff := getParsedDiff(parsedNewSlice, parsedOldSlice, 2)

	fmt.Println(parsedNewSlice)

	fmt.Printf("Difference in Total operations: new results(%s) are differ from old (%s) on %s\n", newSlice[1], oldSlice[1], formattedTOpsDiff)
	fmt.Printf("Difference in ns per operation: new results(%s) are differ from old (%s) on %s\n", newSlice[2], oldSlice[2], formattedNsPerOpsDiff)

	return nil
}

func getParsedDiff(firstSlice []float64, secondSlice []float64, pos int) string {
	diff := firstSlice[pos] / secondSlice[pos] * hundredPercent

	var formattedDiff string
	if diff < hundredPercent {
		diff = hundredPercent - diff
		formattedDiff = fmt.Sprintf("-%f %%", diff)
	} else {
		diff = diff - hundredPercent
		formattedDiff = fmt.Sprintf("+%f %%", diff)
	}

	return formattedDiff
}

func getParsedSlices(firstSlice []string, secondSlice []string) ([]float64, []float64, error) {
	var fs, ss []float64
	for _, e := range firstSlice {
		f, err := strconv.ParseFloat(e, 64)
		if err != nil {
			return []float64{}, []float64{}, err
		}

		fs = append(fs, f)
	}

	for _, e := range secondSlice {
		s, err := strconv.ParseFloat(e, 64)
		if err != nil {
			return []float64{}, []float64{}, err
		}

		ss = append(ss, s)
	}

	return fs, ss, nil
}

func run() error {
	err := initRegexp()
	if err != nil {
		return err
	}

	err = checkArgs()
	if err != nil {
		return err
	}

	oldBenchLine, newBenchLine, err := getBenchlines()
	if err != nil {
		return err
	}

	err = compareResults(oldBenchLine, newBenchLine)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode)
	}
}
