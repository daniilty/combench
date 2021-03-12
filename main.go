package github.com/daniilty/combench

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
		return fmt.Errorf("Usage: ./benchcompare.exe [old.txt] [new.txt]")
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

	tOpsDiff := parsedNewSlice[0] / parsedOldSlice[0] * hundredPercent
	nsPerOpsDiff := parsedNewSlice[1] / parsedOldSlice[1] * hundredPercent

	var formattedTOpsDiff, formattedNsPerOpsDiff string
	if tOpsDiff < hundredPercent {
		tOpsDiff = hundredPercent - tOpsDiff
		formattedTOpsDiff = fmt.Sprintf("-%f %%", tOpsDiff)
	} else {
		tOpsDiff = tOpsDiff - hundredPercent
		formattedTOpsDiff = fmt.Sprintf("+%f %%", tOpsDiff)
	}

	if nsPerOpsDiff < hundredPercent {
		nsPerOpsDiff = hundredPercent - nsPerOpsDiff
		formattedTOpsDiff = fmt.Sprintf("-%f %%", nsPerOpsDiff)
	} else {
		nsPerOpsDiff = nsPerOpsDiff - hundredPercent
		formattedNsPerOpsDiff = fmt.Sprintf("+%f %%", nsPerOpsDiff)
	}

	fmt.Printf("Difference in Total operations: new results(%s) are differ from old (%s) on %s\n", newSlice[1], oldSlice[1], formattedTOpsDiff)
	fmt.Printf("Difference in ns per operation: new results(%s) are differ from old (%s) on %s\n", newSlice[2], oldSlice[2], formattedNsPerOpsDiff)

	return nil
}

func getParsedSlices(firstSlice []string, secondSlice []string) ([]float64, []float64, error) {
	var fs, ss []float64
	for i := 1; i < len(firstSlice); i++ {
		f, err := strconv.ParseFloat(firstSlice[i], 64)
		if err != nil {
			return []float64{}, []float64{}, err
		}

		fs = append(fs, f)
	}

	for i := 1; i < len(secondSlice); i++ {
		s, err := strconv.ParseFloat(secondSlice[i], 64)
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
