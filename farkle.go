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
	// scores contains the scores of all players (including this
	// player).
	scores []int
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

func (c Context) Score(i int) int { return c.scores[i] }

type Strategy interface {
	// Roll represents a roll of the dice.  keep is the dice values that are
	// set aside and not rolled.  Returning all dice (i.e. keep == got) scores
	// all possible points and ends the turn.
	Roll(c Context, got Dice) (keep Dice)
}

type GoForItStrategy int

func (n GoForItStrategy) Roll(c Context, got Dice) (keep Dice) {
	var pts int
	pts, keep = Keep(c.ScoreFn, got)
	if c.Points+pts >= int(n) {
		return got
	}
	return keep
}

type HoldStrategy struct{}

func (_ HoldStrategy) Roll(c Context, got Dice) (keep Dice) {
	var pts int
	pts, keep = Keep(c.ScoreFn, got)
	if c.Points+pts >= c.TurnThresh {
		return got
	}
	return keep
}

func validKeep(got, keep Dice) bool {
	for i, n := range keep {
		if n > got[i] {
			return false
		}
	}
	return true
}

func Turn(ctx Context, rng *rand.Rand, s Strategy) (points int) {
	var d, rem Dice
	var pts int
	n := ndice
	for n > 0 {
		ctx.Points = points
		d = RollDice(rng, n, d)

		// check for failure to roll scoring dice
		if pts, _ := ctx.ScoreFn(0, d); pts == 0 {
			return 0
		}

		keep := s.Roll(ctx, d)
		pts, rem = ctx.ScoreFn(points, keep)
		points += pts
		n -= keep.N()

		if pts == 0 {
			panic("keep dice don't score")
		} else if !validKeep(d, keep) {
			panic("keep dice are not a subset of got dice")
		}

		// check for hot dice
		if n == 0 && rem.N() == 0 {
			n = ndice
		}
	}
	return points
}

type Dice []int

func NewDice() Dice {
	return make(Dice, 7) // num die sides plus one
}

func (d Dice) N() int {
	tot := 0
	for _, n := range d {
		tot += n
	}
	return tot
}

func (d Dice) Clone() Dice {
	clone := make(Dice, len(d))
	for i, n := range d {
		clone[i] = n
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
				scores:     scores,
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
			scores:     scores,
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

func RollDice(rng *rand.Rand, n int, d Dice) Dice {
	if d == nil {
		d = NewDice()
	}
	for i := 0; i < n; i++ {
		d[rng.Intn(6)+1]++
	}
	return d
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
			scoreStraight(prevscore, d.Clone()),
		),
	)
}

func scoreStraight(prevscore int, d Dice) (score int, rem Dice) {
	for _, n := range d {
		if n == 0 {
			return prevscore, d
		}
	}
	return prevscore + 1000, NewDice()
}

func scoreOneFive(prevscore int, d Dice) (score int, rem Dice) {
	d[1] = 0
	d[5] = 0
	return prevscore + d[1]*100 + d[5]*50, d
}

func scoreTriple(prevscore int, d Dice) (score int, rem Dice) {
	for i, n := range d {
		if n >= 3 {
			if i == 1 {
				score += 1000 * (n / 3)
			} else {
				score += i * 100 * (n / 3)
			}
			d[i] = n % 3
		}
	}
	return prevscore + score, d
}
