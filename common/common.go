package common

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Opts struct {
	Part  int    `short:"p" description:"part 1 or 2" required:"true"`
	Input string `short:"i" description:"input file" required:"true"`
}

type Solutions map[int]func(ctx context.Context, log *slog.Logger, opts Opts) error

func ReadAllInput(opts Opts) ([]byte, error) {
	f, err := os.Open(opts.Input)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
