package day3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

// aka `cat day3/input-1.txt | grep -oE 'mul\([0-9]+,[0-9]+\)' | sed -e 's/mul(//' -e 's/)$//' -e 's/,/*/' | paste -sd+ - | bc“
func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	f, err := os.Open(opts.Input)
	if err != nil {
		return err
	}
	defer f.Close()

	cnt, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	getMultsRx := regexp.MustCompile(`mul\([0-9]+,[0-9]+\)`)

	ms := getMultsRx.FindAll(cnt, -1)

	sum := 0
	for _, mult := range ms {
		product, err := parseMult(mult)
		if err != nil {
			return err
		}
		sum += product

	}
	log.Info("result", "sum", sum)

	return nil
}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	// with `do()` and `don't()`
	// do enables stuff and dont disables it

	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	getMultsRx := regexp.MustCompile(`mul\([0-9]+,[0-9]+\)`)
	multOrDoOrDontRx := regexp.MustCompile(fmt.Sprintf(`%s|do\(\)|don't\(\)`, getMultsRx.String()))
	fmt.Printf("multOrDoOrDontRx: %v\n", multOrDoOrDontRx)

	sum := 0
	doing := true
	for pos := 0; pos < len(cnt); pos++ {
		nextThingIdx := multOrDoOrDontRx.FindIndex(cnt[pos:])
		if nextThingIdx == nil {
			break
		}
		pos += nextThingIdx[0]

		nextThing := multOrDoOrDontRx.Find(cnt[pos:])

		if nextThing == nil {
			break
		}
		if string(nextThing) == "do()" {
			doing = true
			continue
		} else if string(nextThing) == "don't()" {
			doing = false
			continue
		}

		if !doing {
			continue
		}

		product, err := parseMult(nextThing)
		if err != nil {
			return err
		}
		sum += product
	}

	log.Info("result", "sum", sum)

	return nil
}

func parseMult(mult []byte) (int, error) {
	getOperandsRx := regexp.MustCompile(`[0-9]+`)

	operands := getOperandsRx.FindAll(mult, -1)
	if len(operands) != 2 {
		return 0, fmt.Errorf("invalid operands: %v", operands)
	}

	left, err := strconv.Atoi(string(operands[0]))
	if err != nil {
		return 0, err
	}

	right, err := strconv.Atoi(string(operands[1]))
	if err != nil {
		return 0, err
	}
	return left * right, nil
}
