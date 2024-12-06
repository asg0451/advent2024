package day5

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"

	"go.coldcutz.net/advent2024/common"
)

var Solutions = common.Solutions{
	1: Part1,
	2: Part2,
}

type rule struct {
	a, b int
}

type ruleset []rule

type update []int

func parseRule(line string) (rule, error) {
	parts := strings.Split(line, "|")
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return rule{}, err
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return rule{}, err
	}
	return rule{a, b}, nil
}

func parseUpdate(line string) (update, error) {
	parts := strings.Split(line, ",")
	u := make(update, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		u = append(u, v)
	}
	return u, nil
}

func parseInput(cnt []byte) (ruleset, []update, error) {
	rules := []rule{}
	updates := []update{}

	ruleSection := true
	for _, line := range strings.Split(string(cnt), "\n") {
		if len(line) == 0 {
			ruleSection = false
			continue
		}
		if ruleSection {
			r, err := parseRule(line)
			if err != nil {
				return nil, nil, err
			}
			rules = append(rules, r)
		} else {
			u, err := parseUpdate(line)
			if err != nil {
				return nil, nil, err
			}
			updates = append(updates, u)
		}
	}

	return rules, updates, nil
}

func (u update) median() int {
	if len(u)%2 == 0 {
		return u[len(u)/2]
	}
	return u[len(u)/2]
}

func (rs ruleset) check(u update) bool {
	rulemap := make(map[int][]int)
	for _, r := range rs {
		rulemap[r.b] = append(rulemap[r.b], r.a)
	}
	forbidden := make(map[int]struct{})
	for _, n := range u {
		if no, ok := rulemap[n]; ok {
			for _, fn := range no {
				forbidden[fn] = struct{}{}
			}
		}
		if _, ok := forbidden[n]; ok {
			return false
		}
	}

	return true
}

func (rs ruleset) checkUpTo(u update, upTo int) bool {
	rulemap := make(map[int][]int)
	for _, r := range rs {
		rulemap[r.b] = append(rulemap[r.b], r.a)
	}
	forbidden := make(map[int]struct{})
	for _, n := range u[:upTo] {
		if no, ok := rulemap[n]; ok {
			for _, fn := range no {
				forbidden[fn] = struct{}{}
			}
		}
		if _, ok := forbidden[n]; ok {
			return false
		}
	}

	return true
}

func Part1(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	rules, updates, err := parseInput(cnt)
	if err != nil {
		return err
	}

	sum := 0
	for _, u := range updates {
		if rules.check(u) {
			log.Debug("valid", "u", u, "median", u.median())
			sum += u.median()
		} else {
			log.Debug("invalid", "u", u)
		}
	}

	log.Info("result", "sum", sum)

	return nil
}

func (rs ruleset) fix(u update, log *slog.Logger) update {
	// lets try swapping stuff until its ok
	// a|b
	// c|d
	// a|c
	// b|c
	// d c a b - [0] bad: swap d and c
	// c d a b - [0] bad: swap with a or b.
	// b d a c - [0] bad: swap with a.
	// a d b c - [0] ok.  [1] bad: swap with c.
	// a c b d - [0] ok.  [1] bad: swap with b.
	// a b c d - [0] ok.  [1] ok., all ok

	// a -> [b, ..]
	rulemap := make(map[int][]int)
	for _, r := range rs {
		rulemap[r.a] = append(rulemap[r.a], r.b)
	}

	for !rs.check(u) {
		log.Debug("trying to fix", "u", u)
		indexed := make(map[int][]int) // map of char to idxs
		for i, n := range u {
			indexed[n] = append(indexed[n], i)
		}

		for i := 0; i < len(u); i++ {
			// is [i] ok?
			if rs.checkUpTo(u, i+1) {
				log.Debug("ok up to", "i", i, "u", u, "sub", u[:i+1])
				continue
			}
			// else swap it with what comes before it
			possibleSwaps := rulemap[u[i]]
			// limit to what's in u
			actuallyPossibleSwaps := make([]int, 0, len(possibleSwaps))
			for i := 0; i < len(possibleSwaps); i++ {
				if _, ok := indexed[possibleSwaps[i]]; ok {
					actuallyPossibleSwaps = append(actuallyPossibleSwaps, possibleSwaps[i])
				}
			}
			if len(actuallyPossibleSwaps) == 0 {
				log.Debug("no possible swaps", "u", u)
				continue
			}

			swapWith := actuallyPossibleSwaps[rand.IntN(len(actuallyPossibleSwaps))]
			swapWithIdxs := indexed[swapWith]
			// pick one
			swapWithIdx := swapWithIdxs[rand.IntN(len(swapWithIdxs))]

			log.Debug("swapping", "ui", u[i], "si", swapWith, "u", u)

			u[i], u[swapWithIdx] = u[swapWithIdx], u[i]
			// update indexed
			indexed = make(map[int][]int)
			for i, n := range u {
				indexed[n] = append(indexed[n], i)
			}
		}
	}

	log.Debug("fixed", "u", u)

	return u
}

func Part2(ctx context.Context, log *slog.Logger, opts common.Opts) error {
	cnt, err := common.ReadAllInput(opts)
	if err != nil {
		return err
	}

	rules, updates, err := parseInput(cnt)
	if err != nil {
		return err
	}

	sum := 0
	for _, u := range updates {
		if rules.check(u) {
			continue
		}
		u := rules.fix(u, log)
		sum += u.median()
	}

	log.Info("result", "sum", sum)

	return nil
}
