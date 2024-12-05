package day4

import (
	"context"
	"log/slog"
	"strings"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	// word search
	const needle = "XMAS"
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	matrix := [][]rune{}
	lines := strings.Split(string(cnt), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		matrix = append(matrix, []rune(line))
	}

	dirs := map[string][2]int{
		"up":         {0, -1},
		"up-left":    {-1, -1},
		"up-right":   {1, -1},
		"down":       {0, 1},
		"down-left":  {-1, 1},
		"down-right": {1, 1},
		"left":       {-1, 0},
		"right":      {1, 0},
	}

	count := 0
	// for each X, search in all directions
	for y, row := range matrix {
		xmax, ymax := len(row)-1, len(matrix)-1
		for x, char := range row {
			if char != rune(needle[0]) {
				continue
			}
			log.Debug("found x at", "x", x, "y", y)
			for dir, dxy := range dirs {
				// try dir
				match := true
				dx, dy := dxy[0], dxy[1]
				for i := 1; i < len(needle); i++ {
					newx, newy := x+dx*i, y+dy*i
					// log.Debug("checking", "dir", dir, "x", newx, "y", newy, "xmax", xmax, "ymax", ymax)
					if newx < 0 || newx > xmax || newy < 0 || newy > ymax {
						match = false
						break
					}
					if matrix[newy][newx] != rune(needle[i]) {
						match = false
						break
					}
				}
				if match {
					log.Debug("found xmas", "dir", dir, "x", x, "y", y)
					count++
				} else {
					log.Debug("giving up on dir", "dir", dir, "x", x, "y", y)
				}
			}
		}
	}

	log.Info("result", "count", count)

	return nil
}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	matrix := [][]rune{}
	lines := strings.Split(string(cnt), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		matrix = append(matrix, []rune(line))
	}

	/// all possible permutations:
	perms := [][][]rune{
		{
			[]rune("M.S"),
			[]rune(".A."),
			[]rune("M.S"),
		},
		{
			[]rune("S.S"),
			[]rune(".A."),
			[]rune("M.M"),
		},
		{
			[]rune("S.M"),
			[]rune(".A."),
			[]rune("S.M"),
		},
		{
			[]rune("M.M"),
			[]rune(".A."),
			[]rune("S.S"),
		},
	}

	// all patterns start with no offset and get incremented as we go
	patterns := []pattern{}
	for _, perm := range perms {
		ms := []matcher{}
		for y, row := range perm {
			for x, r := range row {
				ms = append(ms, matcher{x, y, r})
			}
		}
		patterns = append(patterns, ms)
	}

	count := 0
	for y, row := range matrix {
		for x := range row {
			for pi, pattern := range patterns {
				log.Debug("checking pattern", "pattern", pi, "x", x, "y", y)
				if pattern.match(matrix, x, y) {
					log.Info("found match", "x", x, "y", y, "pattern", pi)
					count++
					break
				}
			}
		}
	}

	log.Info("result", "count", count)

	return nil
}

type matcher struct {
	x, y int
	r    rune
}

func (m matcher) match(matrix [][]rune, x, y int) bool {
	if m.x+x < 0 || m.x+x >= len(matrix[0]) || m.y+y < 0 || m.y+y >= len(matrix) {
		return false
	}
	if m.r == '.' {
		return true
	}
	return matrix[m.y+y][m.x+x] == m.r
}

type pattern []matcher

func (p pattern) match(matrix [][]rune, x, y int) bool {
	for _, m := range p {
		if !m.match(matrix, x, y) {
			return false
		}
	}
	return true
}
