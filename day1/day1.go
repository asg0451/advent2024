package day1

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	left, right, err := readInts(opts)
	if err != nil {
		return err
	}
	sort.Ints(left)
	sort.Ints(right)

	sum := 0
	for i, l := range left {
		r := right[i]
		distance := math.Abs(float64(r - l))
		sum += int(distance)
	}

	log.Info("result", "sum", sum)

	return nil
}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	left, right, err := readInts(opts)
	if err != nil {
		return err
	}

	rightCounts := map[int]int{}
	for _, r := range right {
		rightCounts[r]++
	}

	similarity := 0
	for _, l := range left {
		similarity += l * rightCounts[l]
	}

	log.Info("result", "similarity", similarity)

	return nil
}

func readInts(opts common.Opts) ([]int, []int, error) {
	f, err := os.Open(opts.Input)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	left, right := []int{}, []int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid line: %s", line)
		}
		li, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid line: %s", line)
		}
		ri, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid line: %s", line)
		}
		left = append(left, li)
		right = append(right, ri)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return left, right, nil
}
