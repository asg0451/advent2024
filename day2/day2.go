package day2

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	reports, err := getReports(opts)
	if err != nil {
		return err
	}

	numSafe := 0
	for _, report := range reports {
		if reportIsSafe(report) {
			numSafe++
		}
	}

	log.Info("result", "numSafe", numSafe)

	return nil
}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	reports, err := getReports(opts)
	if err != nil {
		return err
	}

	numSafe := 0
	for _, report := range reports {
		if reportIsSafe(report) {
			numSafe++
			continue
		}

	attempting:
		// try excluding each level in the report
		for i := 0; i < len(report); i++ {
			reportCopy := make([]int, 0, len(report)-1)
			reportCopy = append(reportCopy, report[:i]...)
			reportCopy = append(reportCopy, report[i+1:]...)
			if reportIsSafe(reportCopy) {
				numSafe++
				break attempting
			}
		}
	}

	log.Info("result", "numSafe", numSafe)

	return nil
}

func reportIsSafe(report []int) bool {
	// a report only counts as safe if both of the following are true:
	// - The levels are either all increasing or all decreasing.
	// - Any two adjacent levels differ by at least one and at most three.

	isSafe := true
	isIncreasing := report[0] < report[1]
	for i := 1; i < len(report); i++ {
		isStillIncreasing := report[i-1] < report[i]
		if isIncreasing != isStillIncreasing {
			isSafe = false
			break
		}

		diff := math.Abs(float64(report[i] - report[i-1]))
		if diff < 1 || diff > 3 {
			isSafe = false
			break
		}
	}
	return isSafe
}

func getReports(opts common.Opts) ([][]int, error) {
	f, err := os.Open(opts.Input)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var reports [][]int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		report := make([]int, 0, len(parts))
		for _, levelStr := range parts {
			level, err := strconv.Atoi(levelStr)
			if err != nil {
				return nil, fmt.Errorf("invalid line: %s", line)
			}
			report = append(report, level)
		}
		reports = append(reports, report)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}
