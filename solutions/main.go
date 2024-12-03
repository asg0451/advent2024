package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.coldcutz.net/advent2024/common"
	"go.coldcutz.net/advent2024/day1"
	"go.coldcutz.net/advent2024/day2"
	"go.coldcutz.net/go-stuff/utils"
)

var days = map[int]common.Solutions{
	1: day1.Solutions,
	2: day2.Solutions,
}

type Opts struct {
	Day int `short:"d" description:"day" required:"true"`
	common.Opts
}

func main() {
	ctx, done, log, opts, err := utils.StdSetup[Opts]()
	if err != nil {
		panic(err)
	}
	defer done()

	if err := run(ctx, log, opts); err != nil {
		log.Error("failed to run", "error", err)
		done()
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger, opts Opts) error {
	day, ok := days[opts.Day]
	if !ok {
		return fmt.Errorf("invalid day: %d", opts.Day)
	}

	soln, ok := day[opts.Part]
	if !ok {
		return fmt.Errorf("invalid part: %d", opts.Part)
	}
	return soln(ctx, log, opts.Opts)
}
