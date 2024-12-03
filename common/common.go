package common

import (
	"context"
	"log/slog"
)

type Opts struct {
	Part  int    `short:"p" description:"part 1 or 2" required:"true"`
	Input string `short:"i" description:"input file" required:"true"`
}

type Solutions map[int]func(ctx context.Context, log *slog.Logger, opts Opts) error
