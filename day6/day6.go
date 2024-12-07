package day6

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

type gridEntry rune

const (
	empty      gridEntry = '.'
	obstacle   gridEntry = '#'
	guardUp    gridEntry = '^'
	guardDown  gridEntry = 'v'
	guardLeft  gridEntry = '<'
	guardRight gridEntry = '>'
)

func (ge gridEntry) dirInc(p pos) pos {
	switch ge {
	case guardUp:
		return pos{p[0], p[1] - 1}
	case guardDown:
		return pos{p[0], p[1] + 1}
	case guardLeft:
		return pos{p[0] - 1, p[1]}
	case guardRight:
		return pos{p[0] + 1, p[1]}
	default:
		panic("invalid gridEntry")
	}
}

func (ge gridEntry) turnRight() gridEntry {
	switch ge {
	case guardUp:
		return guardRight
	case guardDown:
		return guardLeft
	case guardLeft:
		return guardUp
	case guardRight:
		return guardDown
	default:
		panic("invalid gridEntry")
	}
}

type grid [][]gridEntry

func (g grid) at(p pos) gridEntry {
	return g[p[1]][p[0]]
}

func (g grid) set(p pos, ge gridEntry) {
	g[p[1]][p[0]] = ge
}

func (g grid) clone() grid {
	clone := make(grid, len(g))
	for i, row := range g {
		clone[i] = make([]gridEntry, len(row))
		copy(clone[i], row)
	}
	return clone
}

func (g grid) String() string {
	var b strings.Builder
	b.WriteRune(' ')
	for i := range g[0] {
		b.WriteRune(rune('0' + i%10))
	}
	b.WriteRune('\n')
	for i, row := range g {
		b.WriteRune(rune('0' + i%10))
		for _, entry := range row {
			b.WriteRune(rune(entry))
		}
		b.WriteRune('\n')
	}
	return b.String()
}

type pos [2]int

func simulateGuard(grid grid, startingPos pos, startingDir gridEntry, log *slog.Logger) int {
	placesVisited := map[pos]struct{}{}

	// guard starts at startingPos
	// guard moves in direction of facing
	// if facing direction is blocked, turn right by 90 degrees
	curPos := startingPos
	curDir := startingDir
	for {
		nextPos := curDir.dirInc(curPos)

		// goes offscreen -- we're done
		if nextPos[0] >= len(grid[0]) || nextPos[1] >= len(grid) || nextPos[0] < 0 || nextPos[1] < 0 {
			return len(placesVisited) + 1 // +1 for the starting position
		}

		nextEntry := grid.at(nextPos)
		if nextEntry == obstacle {
			curDir = curDir.turnRight()
			continue
		}

		// move to next position
		placesVisited[curPos] = struct{}{}
		curPos = nextPos
	}
}

func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	// read into grid
	grid := grid{}
	lines := strings.Split(string(cnt), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, []gridEntry(line))
	}

	fmt.Println(grid)

	// find guard's starting position
	guardPos := pos{-1, -1}
	var guardDir gridEntry
	for y, row := range grid {
		for x, entry := range row {
			switch entry {
			case guardUp, guardDown, guardLeft, guardRight:
				guardPos = pos{x, y}
				guardDir = entry
			}
		}
	}
	if guardPos[0] == -1 {
		return fmt.Errorf("no guard")
	}

	count := simulateGuard(grid, guardPos, guardDir, log)

	log.Info("result", "count", count)

	return nil
}

func simulateGuardStuck(gr grid, startingPos pos, startingDir gridEntry, log *slog.Logger) int {
	// same as simulateGuard, but at each step see if adding an obstacle there gets the guard stuck in a loop
	loopsFound := map[pos]struct{}{}

	placesVisited := map[pos]struct{}{}
	curPos := startingPos
	curDir := startingDir
	for {
		nextPos := curDir.dirInc(curPos)

		// goes offscreen -- we're done
		if nextPos[0] >= len(gr[0]) || nextPos[1] >= len(gr) || nextPos[0] < 0 || nextPos[1] < 0 {
			return len(loopsFound)
		}

		nextEntry := gr.at(nextPos)
		if nextEntry == obstacle {
			curDir = curDir.turnRight()
			continue
		}

		// check if adding an obstacle here would get the guard stuck
		if _, ok := loopsFound[nextPos]; !ok {
			if nextPos != startingPos { // don't add obstacle at starting position
				grClone := gr.clone()
				grClone.set(nextPos, obstacle)
				if doesGuardLoop(grClone, startingPos, startingDir, log) {
					log.Debug("found loop by adding obstacle", "pos", nextPos)
					loopsFound[nextPos] = struct{}{}
				}
			}
		}

		// move to next position
		placesVisited[curPos] = struct{}{}
		curPos = nextPos
	}
}

func doesGuardLoop(grid grid, startingPos pos, startingDir gridEntry, log *slog.Logger) bool {
	// pos -> set of directions we've been in at that pos. if we hit a pos/direction combo we've been in before, we will loop
	visitedDirs := map[pos]map[gridEntry]struct{}{}

	curPos := startingPos
	curDir := startingDir
	for {
		nextPos := curDir.dirInc(curPos)

		// goes offscreen -- we're done
		if nextPos[0] >= len(grid[0]) || nextPos[1] >= len(grid) || nextPos[0] < 0 || nextPos[1] < 0 {
			return false
		}

		nextEntry := grid.at(nextPos)

		if nextEntry == obstacle {
			curDir = curDir.turnRight()
			continue
		}

		// will we loop?
		if prevDirs, ok := visitedDirs[curPos]; ok {
			for prevDir := range prevDirs {
				if prevDir == curDir { // we looped
					return true
				}
			}
		}

		// move to next position
		if _, ok := visitedDirs[curPos]; !ok {
			visitedDirs[curPos] = map[gridEntry]struct{}{}
		}
		visitedDirs[curPos][curDir] = struct{}{}
		curPos = nextPos
	}

}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	// read into grid
	grid := grid{}
	lines := strings.Split(string(cnt), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, []gridEntry(line))
	}

	fmt.Println(grid)

	// find guard's starting position
	guardPos := pos{-1, -1}
	var guardDir gridEntry
	for y, row := range grid {
		for x, entry := range row {
			switch entry {
			case guardUp, guardDown, guardLeft, guardRight:
				guardPos = pos{x, y}
				guardDir = entry
			}
		}
	}
	if guardPos[0] == -1 {
		return fmt.Errorf("no guard")
	}

	count := simulateGuardStuck(grid, guardPos, guardDir, log)

	log.Info("result", "count", count)

	return nil
}
