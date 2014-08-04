package farkle

import (
	"math/rand"
	"time"
)

const (
	towin  = 5000
	thresh = 350
	ndice  = 6
)

type Context struct {
	// Scores contains the scores of all players (including this
	// player).
	Scores []int
	// Index is the player's index - useful for retrieving the player's score.
	Index int
	// Points represents the points accumulated from prior rolls on the
	// currently active turn on - points that are not permanent and may be
	// lose.
	Points int
	// EndScore is the number of points a player must have to end the game.
	EndScore int
	// TurnThresh is the minimum number of points that must be accumulated on a
	// turn for them to become a permanent part of a player's score.
	TurnThresh int
	// ScoreFn is the function used to calculate all scores from the dice.
	ScoreFn ScoreFunc
}

type Strategy interface {
	// Roll represents a roll of the dice.  keep is the dice values that are
	// set aside and not rolled.  Returning zero keep dice (i.e. keep.N() ==
	// 0) ends the turn.
	Roll(c Context, got Dice) (keep Dice)
}

type GoForItStrategy struct{}

func (_ GoForItStrategy) Roll(c Context, got Dice) (keep Dice) {
	_, keep = Keep(c.ScoreFn, got)
	return keep
}

type HoldStrategy struct{}

func (_ HoldStrategy) Roll(c Context, got Dice) (keep Dice) {
	if c.Points >= c.TurnThresh {
		keep = nil
	} else {
		_, keep = Keep(c.ScoreFn, got)
	}
	return keep
}

func ValidKeep(fn ScoreFunc, keep Dice) bool {
	_, rem := fn(0, keep)
	if rem.N() != 0 {
		return false
	}
	return true
}

func Turn(ctx Context, rng *rand.Rand, s Strategy) (points int) {
	n := ndice
	for {
		ctx.Points = points
		got := RollDice(rng, n)

		// check for failure to roll scoring dice
		if pts, _ := ctx.ScoreFn(0, got); pts == 0 {
			return 0
		}

		keep := s.Roll(ctx, got)

		// check for cash-out
		if keep.N() == 0 {
			return points
		} else if !ValidKeep(ctx.ScoreFn, keep) {
			panic("one or more dice set aside are non-scoring")
		}

		points, _ = ctx.ScoreFn(points, keep)
		n -= keep.N()

		// check for hot dice
		if n == 0 {
			n = ndice
		}
	}
	return points
}

type Dice map[int]int

func (d Dice) N() int {
	tot := 0
	for _, n := range d {
		tot += n
	}
	return tot
}

func (d Dice) Clone() Dice {
	clone := Dice{}
	for x, n := range d {
		clone[x] = n
	}
	return clone
}

func Play(rng *rand.Rand, fn ScoreFunc, players ...Strategy) (scores []int) {
	scores = make([]int, len(players))
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().Unix()))
	}
	if fn == nil {
		fn = Score
	}

	for Breaker(scores) < 0 {
		for i, p := range players {
			ctx := Context{
				Scores:     append([]int{}, scores...),
				ScoreFn:    fn,
				Index:      i,
				EndScore:   towin,
				TurnThresh: thresh,
			}
			scores[i] += Turn(ctx, rng, p)
		}
	}

	// give remaining players one more turn
	i := Breaker(scores)
	for i, p := range players[:i] {
		ctx := Context{
			Scores:     append([]int{}, scores...),
			ScoreFn:    fn,
			Index:      i,
			EndScore:   towin,
			TurnThresh: thresh,
		}
		scores[i] += Turn(ctx, rng, p)
	}

	return scores
}

// Breaker returns the index of the first player who broke the game-end
// threshold or -1 if no player has broken it yet.
func Breaker(scores []int) (index int) {
	for i, v := range scores {
		if v > towin {
			return i
		}
	}
	return -1
}

// Winner returns the index of the player with the highest score.
func Winner(scores []int) (index int) {
	best := 0
	for i, v := range scores {
		if v > best {
			best = v
			index = i
		}
	}
	return index
}

func RollDice(rng *rand.Rand, n int) Dice {
	dice := make(Dice, n)
	for i := 0; i < n; i++ {
		dice[rng.Intn(6)+1]++
	}
	return dice
}

func Keep(fn ScoreFunc, d Dice) (points int, scoring Dice) {
	points, rem := fn(0, d)
	scoring = d.Clone()
	for x, n := range rem {
		scoring[x] -= n
	}
	return points, scoring
}

// ScoreFunc returns the highest score possible for the given dice.
type ScoreFunc func(prevscore int, d Dice) (score int, rem Dice)

func Score(prevscore int, d Dice) (score int, rem Dice) {
	return scoreOneFive(
		scoreTriple(
			scoreStraight(prevscore, d),
		),
	)
}

func scoreStraight(prevscore int, d Dice) (score int, rem Dice) {
	for i := 1; i <= 6; i++ {
		if d[i] == 0 {
			return prevscore, d.Clone()
		}
	}
	return prevscore + 1000, Dice{}
}

func scoreOneFive(prevscore int, d Dice) (score int, rem Dice) {
	rem = d.Clone()
	rem[1] = 0
	rem[5] = 0
	return prevscore + d[1]*100 + d[5]*50, rem
}

func scoreTriple(prevscore int, d Dice) (score int, rem Dice) {
	rem = d.Clone()
	for x, n := range d {
		if n >= 3 {
			if x == 1 {
				score += 1000 * (n / 3)
			} else {
				score += x * 100 * (n / 3)
			}
			rem[x] = n % 3
		}
	}
	return prevscore + score, rem
}
